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
	PayOrder(*proto.CreatePaymentRequest) (*proto.OrderPayment, error)
	Transaction(*proto.PaymentTransaction) (*proto.EmptyPayment, error)
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

func (u *PaymentRepositoryImpl) PayOrder(payload *proto.CreatePaymentRequest) (*proto.OrderPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.PayOrder(ctx, payload)
}

func (u *PaymentRepositoryImpl) Transaction(payload *proto.PaymentTransaction) (*proto.EmptyPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.Transaction(ctx, payload)
}
