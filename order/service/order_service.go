package service

import (
	"context"
	"database/sql"
	"errors"
	"order/cmd/db"
	"order/helper"
	"order/proto"
	"order/repository"
	"order/transport/kafka"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	orderCacheTTL     = 10 * time.Minute
	orderListCacheTTL = 5 * time.Minute
)

type OrderService struct {
	DB            *sql.DB
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemsRepository
	productRepo   repository.ProductRepository
	ctx           context.Context
}

func NewOrderItemService(DB *sql.DB, orderRepo repository.OrderRepository, orderItemRepo repository.OrderItemsRepository, productRepo repository.ProductRepository, ctx context.Context) *OrderService {
	return &OrderService{
		DB:            DB,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		productRepo:   productRepo,
		ctx:           ctx,
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

	if err := db.DeleteCacheByPattern(u.ctx, "orders:user*"); err != nil {
		logrus.Warnf("Failed to invalidate product list cache: %v", err)
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
	page := (payload.Offset / 15) + 1
	key := db.RedisOrderKey(int(payload.UserId), int(page))

	cachedList, err := db.GetCacheOrderList(u.ctx, key)
	if err == nil {
		logrus.Infof("Cache HIT for order list page: %d", page)
		return cachedList, nil
	}
	logrus.Infof("Cache MISS for order list page: %d", page)

	orderResponse, err := u.orderRepo.GetOrderByUserID(payload, u.DB)
	if err != nil {
		return nil, err
	}

	if err := db.SetCache(u.ctx, key, orderResponse, orderListCacheTTL); err != nil {
		logrus.Warnf("Failed to cache product list: %v", err)
	} else {
		logrus.Debugf("Product list page %d cached successfully", page)
	}
	return orderResponse, nil
}

func (u *OrderService) GetOrderById(payload *proto.GetOrderByIdRequest) (*proto.Order, error) {
	// Try to get from cache first
	key := db.RedisOrderByIdKey(int(payload.OrderId))

	cachedOrder, err := db.GetCacheOrder(u.ctx, key)
	if err == nil && cachedOrder != nil {
		logrus.Infof("Cache HIT for order ID: %d", payload.OrderId)
		return cachedOrder, nil
	}
	logrus.Infof("Cache MISS for order ID: %d", payload.OrderId)

	// Get from database
	order, err := u.orderRepo.GetOrderById(int(payload.OrderId), int(payload.UserId), u.DB)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// Cache the result
	if err := db.SetCache(u.ctx, key, order, orderCacheTTL); err != nil {
		logrus.Warnf("Failed to cache order: %v", err)
	} else {
		logrus.Debugf("Order ID %d cached successfully", payload.OrderId)
	}

	return order, nil
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

	// Invalidate order list cache (all users)
	if err := db.DeleteCacheByPattern(u.ctx, "orders:user*"); err != nil {
		logrus.Warnf("Failed to invalidate order list cache: %v", err)
	}

	// Invalidate single order cache
	orderCacheKey := db.RedisOrderByIdKey(orderID)
	if err := db.DeleteCache(u.ctx, orderCacheKey); err != nil {
		logrus.Warnf("Failed to invalidate order cache for ID %d: %v", orderID, err)
	} else {
		logrus.Infof("Successfully invalidated cache for order ID: %d", orderID)
	}

	return nil
}
