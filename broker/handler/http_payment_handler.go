package handler

import (
	"broker/middleware"
	"broker/proto"
	"broker/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	repo repository.PaymentRepository
}

func NewPaymentHandler(repo repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		repo: repo,
	}
}

func (u *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	// Protected routes (require auth)
	paymentRoutes := r.Group("/payment")
	paymentRoutes.Use(middleware.ProtectedEndpoint())
	paymentRoutes.GET("/order/:order_id", u.GetPaymentByOrderId)
	paymentRoutes.POST("/initiate", u.InitiatePayment)

	// Webhook route (no auth - called by Midtrans)
	r.POST("/payment/webhook/midtrans", u.HandleMidtransWebhook)
}

// GetPaymentByOrderId retrieves payment by order ID
func (u *PaymentHandler) GetPaymentByOrderId(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := u.repo.GetPaymentByOrderId(int32(orderID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, payment)
}

// InitiatePayment creates a Midtrans Snap transaction
func (u *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req struct {
		OrderID        int32  `json:"order_id" binding:"required"`
		PaymentMethod  string `json:"payment_method"`
		PaymentChannel string `json:"payment_channel"`
		CustomerName   string `json:"customer_name" binding:"required"`
		CustomerEmail  string `json:"customer_email" binding:"required"`
		CustomerPhone  string `json:"customer_phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("Initiating payment for order: %d", req.OrderID)

	response, err := u.repo.InitiatePayment(&proto.InitiatePaymentRequest{
		OrderId:        req.OrderID,
		PaymentMethod:  req.PaymentMethod,
		PaymentChannel: req.PaymentChannel,
		CustomerName:   req.CustomerName,
		CustomerEmail:  req.CustomerEmail,
		CustomerPhone:  req.CustomerPhone,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"payment_id":   response.PaymentId,
		"token":        response.GatewayToken,
		"redirect_url": response.GatewayRedirectUrl,
		"va_number":    response.VaNumber,
		"qr_code_url":  response.QrCodeUrl,
		"expired_at":   response.ExpiredAt,
		"status":       response.Status,
	})
}

// HandleMidtransWebhook processes webhook notifications from Midtrans
func (u *PaymentHandler) HandleMidtransWebhook(c *gin.Context) {
	var req struct {
		OrderID           string `json:"order_id"`
		TransactionID     string `json:"transaction_id"`
		TransactionStatus string `json:"transaction_status"`
		PaymentType       string `json:"payment_type"`
		GrossAmount       string `json:"gross_amount"`
		SignatureKey      string `json:"signature_key"`
		FraudStatus       string `json:"fraud_status"`
		StatusCode        string `json:"status_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Invalid webhook payload: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("Received Midtrans webhook for order: %s, status: %s", req.OrderID, req.TransactionStatus)

	err := u.repo.HandleWebhook(&proto.WebhookRequest{
		OrderId:           req.OrderID,
		TransactionId:     req.TransactionID,
		TransactionStatus: req.TransactionStatus,
		PaymentType:       req.PaymentType,
		GrossAmount:       req.GrossAmount,
		SignatureKey:      req.SignatureKey,
		FraudStatus:       req.FraudStatus,
		StatusCode:        req.StatusCode,
	})
	if err != nil {
		logrus.Errorf("Failed to handle webhook: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
