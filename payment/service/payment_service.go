package service

import (
	"context"
	"database/sql"
	"errors"
	"payment/helper"
	"payment/proto"
	"payment/repository"
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

func (u *PaymentService) AddPayment(payment *proto.CreatePaymentRequest) (*proto.OrderPayment, error) {
	logrus.Info("Begin transaction")
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("create payment")
	paymentID, err := u.paymentRepo.CreatePayment(payment, tx)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("error when create payment: %v", err)
		return nil, err
	}

	return &proto.OrderPayment{
		Id:         int32(paymentID),
		OrderId:    payment.OrderId,
		UserId:     payment.UserId,
		Status:     "pending",
		TotalPrice: payment.TotalPrice,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (u *PaymentService) Transaction(transaction *proto.PaymentTransaction) error {
	logrus.Info("create transaction")
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("get payment by id")
	payment, err := u.paymentRepo.GetByID(u.ctx, int(transaction.PaymentId), tx)
	if err != nil {
		return err
	}

	logrus.Info("check payment if already paid")
	if payment.Status == "paid" {
		return errors.New("payment already paid")
	}

	logrus.Info("check payment")
	if transaction.Money < payment.TotalPrice {
		if _, err := u.orderRepo.UpdateOrderStatus(u.ctx, &proto.UpdateOrderStatusRequest{OrderId: payment.OrderId, Status: "failed"}); err != nil {
			logrus.Errorf("error when update order status: %v", err)
			return err
		}
		return errors.New("money not enough")
	}

	logrus.Info("update payment")
	if err := u.paymentRepo.UpdatePayment(u.ctx, "paid", int(transaction.PaymentId), tx); err != nil {
		if _, err := u.orderRepo.UpdateOrderStatus(u.ctx, &proto.UpdateOrderStatusRequest{OrderId: payment.OrderId, Status: "failed"}); err != nil {
			return err
		}
		tx.Rollback()
		return err
	}

	logrus.Info("update order status")
	if _, err := u.orderRepo.UpdateOrderStatus(u.ctx, &proto.UpdateOrderStatusRequest{OrderId: payment.OrderId, Status: "paid"}); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (u *PaymentService) GetPaymentByOrderId(orderID int) (*proto.OrderPayment, error) {
	logrus.Infof("get payment by order id: %d", orderID)
	return u.paymentRepo.GetByOrderID(u.ctx, orderID, u.DB)
}
