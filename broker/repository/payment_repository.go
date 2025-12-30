package repository

import (
	"broker/proto"
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentRepository interface {
	GetPaymentByOrderId(orderID int32) (*proto.PaymentResponse, error)
	InitiatePayment(req *proto.InitiatePaymentRequest) (*proto.InitiatePaymentResponse, error)
	HandleWebhook(req *proto.WebhookRequest) error
}

type PaymentRepositoryImpl struct {
	client proto.PaymentServiceClient
}

func NewPaymentRepository() *PaymentRepositoryImpl {
	addr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if addr == "" {
		addr = "payment-service:60001"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect to payment service: %v", err)
	}

	logrus.Infof("Connected to payment service at %s", addr)
	client := proto.NewPaymentServiceClient(conn)

	return &PaymentRepositoryImpl{client: client}
}

func (u *PaymentRepositoryImpl) GetPaymentByOrderId(orderID int32) (*proto.PaymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.GetPaymentByOrderId(ctx, &proto.GetPaymentByOrderIdRequest{OrderId: orderID})
}

func (u *PaymentRepositoryImpl) InitiatePayment(req *proto.InitiatePaymentRequest) (*proto.InitiatePaymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return u.client.InitiatePayment(ctx, req)
}

func (u *PaymentRepositoryImpl) HandleWebhook(req *proto.WebhookRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := u.client.HandleWebhook(ctx, req)
	return err
}
