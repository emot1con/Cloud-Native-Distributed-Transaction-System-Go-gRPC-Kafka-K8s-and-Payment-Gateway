package grpc

import (
	"context"
	"net"
	"order/cmd/db"
	"order/proto"
	"order/repository"
	"order/service"
	"order/transport/kafka"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type OrderGRPCServer struct {
	service *service.OrderService
	proto.UnimplementedOrderServiceServer
}

func NewOrderGRPCServer(service *service.OrderService) *OrderGRPCServer {
	return &OrderGRPCServer{
		service: service,
	}
}

func (u *OrderGRPCServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	order, err := u.service.CreateOrder(req)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderGRPCServer) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrdersResponse, error) {
	orders, err := u.service.GetOrderByUserID(req)
	if err != nil {
		return nil, err
	}

	return &proto.OrdersResponse{
		Orders: orders,
	}, nil
}

func (u *OrderGRPCServer) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.EmptyOrder, error) {
	if err := u.service.UpdateOrderStatus(req.Status, int(req.OrderId)); err != nil {
		return nil, err
	}

	return &proto.EmptyOrder{}, nil
}

func GRPCListen() {
	DB, err := db.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	addr := []string{os.Getenv("KAFKA_BROKER_URL")}

	ctx := context.Background()

	orderRepo := repository.NewOrderRepositoryImpl()
	orderItemRepo := repository.NewOrderItemsRepositoryImpl()
	productRepo := repository.NewProductRepositoryImpl()

	orderService := service.NewOrderItemService(DB, orderRepo, orderItemRepo, productRepo, ctx)
	orderGRPC := NewOrderGRPCServer(orderService)

	if err := kafka.ConnectProducer(addr); err != nil {
		logrus.Fatalf("failed to connect to kafka: %v", err)
	}

	conn, err := net.Listen("tcp", ":30001")
	if err != nil {
		logrus.Fatalf("failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterOrderServiceServer(srv, orderGRPC)
	logrus.Infof("gRPC Server started on port 30001")

	if err := srv.Serve(conn); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}
