package service

import (
	"database/sql"
	"errors"
	"order/helper"
	"order/proto"
	"order/repository"
	"order/transport/kafka"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderService struct {
	DB            *sql.DB
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemsRepository
	productRepo   repository.ProductRepository
}

func NewOrderItemService(DB *sql.DB, orderRepo repository.OrderRepository, orderItemRepo repository.OrderItemsRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		DB:            DB,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		productRepo:   productRepo,
	}
}

func (u *OrderService) CreateOrder(payload *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	var totalPrices float64

	addr := []string{os.Getenv("KAFKA_BROKER_URL")}
	topic := os.Getenv("KAFKA_ORDER_TOPIC")

	var products []*proto.Product

	logrus.Info("Calculating total price")
	for _, v := range payload.Items {
		product, err := u.productRepo.GetProduct(&proto.GetProductRequest{Id: v.ProductId})
		if err != nil {
			return nil, err
		}
		products = append(products, product)

		if product.Stock < v.Quantity {
			return nil, errors.New("stock is not enough")
		}
		totalPrice := float64(v.Quantity) * product.Price
		totalPrices += totalPrice
	}
	payload.TotalPrice = totalPrices

	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}

	rollback := true
	defer func() {
		if rollback {
			if rErr := tx.Rollback(); rErr != nil {
				logrus.Errorf("Rollback error: %v", rErr)
			}
		}
	}()

	logrus.Info("Create order")
	orderID, err := u.orderRepo.CreateOrder(payload, payload.TotalPrice, tx)
	if err != nil {
		return nil, err
	}

	logrus.Info("Saving order items")
	for _, v := range payload.Items {
		orderItemPayload := &proto.OrderItemRequest{
			OrderId:   int32(orderID),
			ProductId: v.ProductId,
			Quantity:  v.Quantity,
		}
		err := u.orderItemRepo.CreateOrderItems(orderItemPayload, tx)
		if err != nil {
			return nil, err
		}
	}

	logrus.Info("Updating products")
	for i, v := range payload.Items {
		if _, err := u.productRepo.UpdateProduct(&proto.Product{
			Id:          products[i].Id,
			Name:        products[i].Name,
			Description: products[i].Description,
			Price:       products[i].Price,
			Stock:       products[i].Stock - v.Quantity,
			CreatedAt:   products[i].CreatedAt,
			UpdatedAt:   products[i].UpdatedAt,
		}); err != nil {
			return nil, err
		}
	}

	// Commit transaction sebelum publish Kafka
	logrus.Info("Committing transaction")
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	rollback = false

	logrus.Info("Sending message to kafka")
	partition, offset, err := kafka.SendMessage(addr, topic, &proto.Order{
		Id:         int32(orderID),
		UserId:     payload.UserId,
		Status:     "Pending",
		TotalPrice: payload.TotalPrice,
	})
	if err != nil {
		logrus.Errorf("Failed to send message to Kafka: %v", err)
	} else {
		logrus.Infof("Message sent to topic: %s partition: %d, offset: %d", topic, partition, offset)
	}

	return &proto.OrderResponse{
		Order: &proto.Order{
			Id:         int32(orderID),
			TotalPrice: payload.TotalPrice,
			CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (u *OrderService) GetOrderByUserID(payload *proto.GetOrderRequest) ([]*proto.Order, error) {
	orderResponse, err := u.orderRepo.GetOrderByUserID(payload, u.DB)
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

func (u *OrderService) UpdateOrderStatus(status string, orderID int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	if err := u.orderRepo.UpdateOrderStatus(status, orderID, tx); err != nil {
		return err
	}
	return nil
}
