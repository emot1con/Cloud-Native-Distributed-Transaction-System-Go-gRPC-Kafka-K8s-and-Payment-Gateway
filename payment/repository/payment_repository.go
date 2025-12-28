package repository

import (
	"context"
	"database/sql"
	"errors"
	"payment/proto"
	"time"
)

type PaymentRepository interface {
	CreatePayment(payload *proto.CreatePaymentRequest, tx *sql.Tx) (int, error)
	UpdatePayment(ctx context.Context, status string, ID int, tx *sql.Tx) error
	GetByID(ctx context.Context, ID int, tx *sql.Tx) (*proto.OrderPayment, error)
	GetByOrderID(ctx context.Context, orderID int, db *sql.DB) (*proto.OrderPayment, error)
	DeletePayment(ctx context.Context, ID int, tx *sql.Tx) error
}

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{}
}

func (u *PaymentRepositoryImpl) CreatePayment(payload *proto.CreatePaymentRequest, tx *sql.Tx) (int, error) {
	var paymentID int
	SQL := "INSERT INTO payments(order_id, user_id, total_price) VALUES ($1, $2, $3) returning id"
	if err := tx.QueryRow(SQL, payload.OrderId, payload.UserId, payload.TotalPrice).Scan(&paymentID); err != nil {
		return 0, err
	}

	return paymentID, nil
}

func (u *PaymentRepositoryImpl) UpdatePayment(ctx context.Context, status string, ID int, tx *sql.Tx) error {
	loc := time.FixedZone("WIB", 7*60*60)
	SQL := `UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3`
	now := time.Now().In(loc)
	if _, err := tx.ExecContext(ctx, SQL, status, now, ID); err != nil {
		return err
	}

	return nil
}

func (u PaymentRepositoryImpl) GetByID(ctx context.Context, ID int, tx *sql.Tx) (*proto.OrderPayment, error) {
	SQL := "SELECT id, order_id, user_id, status, total_price, created_at, updated_at FROM payments WHERE id = $1"
	rows := tx.QueryRowContext(ctx, SQL, ID)

	orderPayment := &proto.OrderPayment{}
	if err := rows.Scan(
		&orderPayment.Id,
		&orderPayment.OrderId,
		&orderPayment.UserId,
		&orderPayment.Status,
		&orderPayment.TotalPrice,
		&orderPayment.CreatedAt,
		&orderPayment.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order payment not found")
		}
		return nil, err
	}
	return orderPayment, nil
}

func (u *PaymentRepositoryImpl) DeletePayment(ctx context.Context, ID int, tx *sql.Tx) error {
	SQL := "DELETE FROM order_items WHERE id = $1"
	if _, err := tx.Exec(SQL, ID); err != nil {
		return err
	}

	return nil
}

func (u *PaymentRepositoryImpl) GetByOrderID(ctx context.Context, orderID int, db *sql.DB) (*proto.OrderPayment, error) {
	SQL := "SELECT id, order_id, user_id, status, total_price, created_at, updated_at FROM payments WHERE order_id = $1"
	row := db.QueryRowContext(ctx, SQL, orderID)

	orderPayment := &proto.OrderPayment{}
	if err := row.Scan(
		&orderPayment.Id,
		&orderPayment.OrderId,
		&orderPayment.UserId,
		&orderPayment.Status,
		&orderPayment.TotalPrice,
		&orderPayment.CreatedAt,
		&orderPayment.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found for this order")
		}
		return nil, err
	}
	return orderPayment, nil
}
