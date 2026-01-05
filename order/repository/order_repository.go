package repository

import (
	"database/sql"
	"order/proto"
	"time"
)

type OrderRepository interface {
	CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error)
	GetOrderByUserID(payload *proto.GetOrderRequest, db *sql.DB) ([]*proto.Order, error)
	GetOrderById(orderID int, userID int, db *sql.DB) (*proto.Order, error)
	UpdateOrderStatus(status string, orderID int, tx *sql.Tx) error
}

type OrderRepositoryImpl struct{}

func NewOrderRepositoryImpl() *OrderRepositoryImpl {
	return &OrderRepositoryImpl{}
}

func (u *OrderRepositoryImpl) CreateOrder(payload *proto.CreateOrderRequest, price float64, tx *sql.Tx) (int, error) {
	SQL := "INSERT INTO orders(user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id"
	var orderID int
	err := tx.QueryRow(SQL, payload.UserId, price, "pending").Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *OrderRepositoryImpl) GetOrderByUserID(payload *proto.GetOrderRequest, db *sql.DB) ([]*proto.Order, error) {
	SQL := "SELECT id, user_id, status, total_price, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY id ASC LIMIT 15 OFFSET $2"
	rows, err := db.Query(SQL, payload.UserId, payload.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*proto.Order
	for rows.Next() {
		order := &proto.Order{}
		if err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (u *OrderRepositoryImpl) GetOrderById(orderID int, userID int, db *sql.DB) (*proto.Order, error) {
	SQL := "SELECT id, user_id, status, total_price, created_at, updated_at FROM orders WHERE id = $1 AND user_id = $2"
	row := db.QueryRow(SQL, orderID, userID)

	order := &proto.Order{}
	if err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Status,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order not found
		}
		return nil, err
	}

	return order, nil
}

func (u *OrderRepositoryImpl) UpdateOrderStatus(status string, orderID int, tx *sql.Tx) error {
	loc := time.FixedZone("WIB", 7*60*60)
	SQL := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	now := time.Now().In(loc)
	if _, err := tx.Exec(SQL, status, now, orderID); err != nil {
		return err
	}
	return nil
}
