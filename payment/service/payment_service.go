package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment/client"
	"payment/proto"
	"payment/repository"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	DB          *sql.DB
	ctx         context.Context
}

func NewPaymentService(repo repository.PaymentRepository, DB *sql.DB, ctx context.Context, orderRepo repository.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: repo,
		orderRepo:   orderRepo,
		DB:          DB,
		ctx:         ctx,
	}
}

// CreatePayment creates a new payment record when order is created (via Kafka)
func (u *PaymentService) CreatePayment(payment *proto.CreatePaymentRequest) (*proto.PaymentResponse, error) {
	logrus.Infof("Creating payment for order: %d, amount: %f", payment.OrderId, payment.Amount)

	paymentResponse, err := u.paymentRepo.CreatePayment(u.ctx, payment, u.DB)
	if err != nil {
		logrus.Errorf("error when create payment: %v", err)
		return nil, err
	}

	logrus.Infof("Payment created with ID: %d", paymentResponse.Id)
	return paymentResponse, nil
}

// GetPaymentByOrderId retrieves payment by order ID
func (u *PaymentService) GetPaymentByOrderId(orderID int) (*proto.PaymentResponse, error) {
	logrus.Infof("Get payment by order id: %d", orderID)
	return u.paymentRepo.GetByOrderID(u.ctx, orderID, u.DB)
}

// InitiatePayment creates a Midtrans Snap transaction
func (u *PaymentService) InitiatePayment(req *proto.InitiatePaymentRequest) (*proto.InitiatePaymentResponse, error) {
	logrus.Infof("Initiating payment for order: %d", req.OrderId)

	// Get existing payment
	payment, err := u.paymentRepo.GetByOrderID(u.ctx, int(req.OrderId), u.DB)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}

	// Check if already paid
	if payment.Status == "paid" || payment.Status == "success" {
		return nil, errors.New("payment already completed")
	}

	// Check if already has gateway token (already initiated)
	if payment.GatewayToken != "" && payment.Status == "pending" {
		// Return existing token
		return &proto.InitiatePaymentResponse{
			PaymentId:          payment.Id,
			GatewayToken:       payment.GatewayToken,
			GatewayRedirectUrl: payment.GatewayRedirectUrl,
			VaNumber:           payment.VaNumber,
			QrCodeUrl:          payment.QrCodeUrl,
			ExpiredAt:          payment.ExpiredAt,
			Status:             payment.Status,
		}, nil
	}

	// Generate unique order ID for Midtrans
	gatewayOrderID := client.GenerateOrderID(payment.Id)

	// Create Snap transaction
	snapResp, err := client.CreateSnapTransaction(
		gatewayOrderID,
		int64(payment.Amount),
		req.CustomerName,
		req.CustomerEmail,
		req.CustomerPhone,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create snap transaction: %v", err)
	}

	// Calculate expiry time (24 hours from now)
	expiredAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	// Update payment with gateway info
	payment.PaymentMethod = req.PaymentMethod
	payment.PaymentChannel = req.PaymentChannel
	payment.GatewayOrderId = gatewayOrderID
	payment.GatewayToken = snapResp.Token
	payment.GatewayRedirectUrl = snapResp.RedirectURL
	payment.ExpiredAt = expiredAt
	payment.Status = "pending"

	if err := u.paymentRepo.UpdatePaymentGateway(u.ctx, payment, u.DB); err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	return &proto.InitiatePaymentResponse{
		PaymentId:          payment.Id,
		GatewayToken:       snapResp.Token,
		GatewayRedirectUrl: snapResp.RedirectURL,
		ExpiredAt:          expiredAt,
		Status:             "pending",
	}, nil
}

// HandleWebhook processes webhook notifications from Midtrans
func (u *PaymentService) HandleWebhook(req *proto.WebhookRequest) error {
	logrus.Infof("Handling webhook for order: %s, status: %s", req.OrderId, req.TransactionStatus)

	// Verify signature
	if !client.VerifySignature(req.OrderId, req.StatusCode, req.GrossAmount, req.SignatureKey) {
		return errors.New("invalid webhook signature")
	}

	// Map Midtrans status to our status
	status := client.MapTransactionStatus(req.TransactionStatus, req.FraudStatus)

	logrus.Infof("Mapped status: %s -> %s", req.TransactionStatus, status)

	// Check if this is a Midtrans test notification
	if strings.HasPrefix(req.OrderId, "payment_notif_test_") {
		logrus.Info("Received Midtrans test notification, returning success")
		return nil
	}

	// Parse order ID to get payment ID (format: PAY-{payment_id}-{timestamp})
	var paymentID int
	var timestamp int64
	_, err := fmt.Sscanf(req.OrderId, "PAY-%d-%d", &paymentID, &timestamp)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %v", err)
	}

	// Get payment by payment_id (not order_id - the parsed paymentID is the payment.id)
	payment, err := u.paymentRepo.GetByID(u.ctx, paymentID, u.DB)
	if err != nil {
		// Try to find by gateway_order_id instead
		logrus.Warnf("Payment not found by ID %d, searching by gateway_order_id: %s", paymentID, req.OrderId)
		payment, err = u.paymentRepo.GetByGatewayOrderID(u.ctx, req.OrderId, u.DB)
		if err != nil {
			return fmt.Errorf("payment not found: %v", err)
		}
	}

	if payment == nil {
		return fmt.Errorf("payment not found for order ID: %s", req.OrderId)
	}

	// Update payment status
	if err := u.paymentRepo.UpdatePaymentStatus(u.ctx, int(payment.OrderId), status, req.TransactionId, u.DB); err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// If payment successful, update order status via gRPC
	if status == "paid" {
		logrus.Infof("Payment successful, updating order status for order: %d", payment.OrderId)
		if _, err := u.orderRepo.UpdateOrderStatus(u.ctx, &proto.UpdateOrderStatusRequest{
			OrderId: payment.OrderId,
			Status:  "paid",
		}); err != nil {
			logrus.Errorf("Failed to update order status: %v", err)
			return err
		}
	} else if status == "failed" {
		logrus.Infof("Payment failed, updating order status for order: %d", payment.OrderId)
		if _, err := u.orderRepo.UpdateOrderStatus(u.ctx, &proto.UpdateOrderStatusRequest{
			OrderId: payment.OrderId,
			Status:  "failed",
		}); err != nil {
			logrus.Errorf("Failed to update order status: %v", err)
		}
	}

	return nil
}
