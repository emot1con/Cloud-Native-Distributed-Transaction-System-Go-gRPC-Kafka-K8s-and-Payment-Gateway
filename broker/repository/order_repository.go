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

type OrderRepository interface {
	CreateOrder(*proto.CreateOrderRequest) (*proto.OrderResponse, error)
	GetOrder(*proto.GetOrderRequest) (*proto.OrdersResponse, error)
	GetOrderById(*proto.GetOrderByIdRequest) (*proto.OrderResponse, error)
}

type OrderRepositoryImpl struct {
	client proto.OrderServiceClient
}

func NewOrderRepository() *OrderRepositoryImpl {
	addr := os.Getenv("ORDER_SERVICE_ADDR")
	if addr == "" {
		addr = "order-service:30001"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect to order service: %v", err)
	}

	logrus.Infof("Connected to order service at %s", addr)
	client := proto.NewOrderServiceClient(conn)

	return &OrderRepositoryImpl{client: client}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Create the order and pass to order service with")
	return u.client.CreateOrder(ctx, payload)
}

func (u *OrderRepositoryImpl) GetOrder(payload *proto.GetOrderRequest) (*proto.OrdersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Get the order by ID and pass to order service with user ID: %d", payload.UserId)
	return u.client.GetOrder(ctx, payload)
}

func (u *OrderRepositoryImpl) GetOrderById(payload *proto.GetOrderByIdRequest) (*proto.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Get order by ID: %d for user ID: %d", payload.OrderId, payload.UserId)
	return u.client.GetOrderById(ctx, payload)
}
