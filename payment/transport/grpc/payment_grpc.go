package grpc

import (
	"context"
	"net"
	"payment/client"
	"payment/cmd/db"
	"payment/proto"
	"payment/repository"
	"payment/service"
	"payment/transport/kafka"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type PaymentGRPCServer struct {
	service *service.PaymentService
	proto.UnimplementedPaymentServiceServer
}

func NewPaymentGRPCServer(service *service.PaymentService) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		service: service,
	}
}

// CreatePayment creates a new payment record (called when order is created)
func (u *PaymentGRPCServer) CreatePayment(ctx context.Context, req *proto.CreatePaymentRequest) (*proto.PaymentResponse, error) {
	payment, err := u.service.CreatePayment(req)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentByOrderId retrieves payment by order ID
func (u *PaymentGRPCServer) GetPaymentByOrderId(ctx context.Context, req *proto.GetPaymentByOrderIdRequest) (*proto.PaymentResponse, error) {
	payment, err := u.service.GetPaymentByOrderId(int(req.OrderId))
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// InitiatePayment creates a Midtrans Snap transaction
func (u *PaymentGRPCServer) InitiatePayment(ctx context.Context, req *proto.InitiatePaymentRequest) (*proto.InitiatePaymentResponse, error) {
	response, err := u.service.InitiatePayment(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// HandleWebhook processes webhook notifications from Midtrans
func (u *PaymentGRPCServer) HandleWebhook(ctx context.Context, req *proto.WebhookRequest) (*proto.EmptyPayment, error) {
	if err := u.service.HandleWebhook(req); err != nil {
		return nil, err
	}
	return &proto.EmptyPayment{}, nil
}

func GRPCListen(addr []string, topic []string, groupID string) {
	// Initialize Midtrans client
	client.InitMidtransClient()

	DB, err := db.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	ctx := context.Background()
	paymentRepo := repository.NewPaymentRepository()
	orderRepo := repository.NewOrderRepository()
	service := service.NewPaymentService(paymentRepo, DB, ctx, orderRepo)
	connection := NewPaymentGRPCServer(service)

	lis, err := net.Listen("tcp", ":60001")
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterPaymentServiceServer(srv, connection)
	logrus.Infof("gRPC Server started on port 60001")

	go func() {
		if err := srv.Serve(lis); err != nil {
			logrus.Fatalf("error when connect to gRPC Server: %v", err)
		}
	}()
	go kafka.ProcessMessage(addr, topic, groupID, service)
}
