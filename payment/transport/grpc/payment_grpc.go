package grpc

import (
	"context"
	"net"
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

func (u *PaymentGRPCServer) PayOrder(ctx context.Context, req *proto.CreatePaymentRequest) (*proto.OrderPayment, error) {
	payment, err := u.service.AddPayment(req)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *PaymentGRPCServer) Transaction(ctx context.Context, req *proto.PaymentTransaction) (*proto.EmptyPayment, error) {
	if err := u.service.Transaction(req); err != nil {
		return nil, err
	}

	return &proto.EmptyPayment{}, nil
}

func (u *PaymentGRPCServer) GetPaymentByOrderId(ctx context.Context, req *proto.GetPaymentByOrderIdRequest) (*proto.OrderPayment, error) {
	payment, err := u.service.GetPaymentByOrderId(int(req.OrderId))
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func GRPCListen(addr []string, topic []string, groupID string) {
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
