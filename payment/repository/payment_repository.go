package repository

import (
	"context"
	"database/sql"
	"errors"
	"payment/proto"
	"time"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payload *proto.CreatePaymentRequest, db *sql.DB) (*proto.PaymentResponse, error)
	GetByID(ctx context.Context, paymentID int, db *sql.DB) (*proto.PaymentResponse, error)
	GetByOrderID(ctx context.Context, orderID int, db *sql.DB) (*proto.PaymentResponse, error)
	GetByGatewayOrderID(ctx context.Context, gatewayOrderID string, db *sql.DB) (*proto.PaymentResponse, error)
	UpdatePaymentGateway(ctx context.Context, payment *proto.PaymentResponse, db *sql.DB) error
	UpdatePaymentStatus(ctx context.Context, orderID int, status string, transactionID string, db *sql.DB) error
}

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{}
}

func (u *PaymentRepositoryImpl) CreatePayment(ctx context.Context, payload *proto.CreatePaymentRequest, db *sql.DB) (*proto.PaymentResponse, error) {
	SQL := `INSERT INTO payments(order_id, amount, status) 
			VALUES ($1, $2, 'pending') 
			RETURNING id, order_id, amount, currency, status, created_at`

	row := db.QueryRowContext(ctx, SQL, payload.OrderId, payload.Amount)

	payment := &proto.PaymentResponse{}
	var createdAt time.Time
	var currency sql.NullString

	if err := row.Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Amount,
		&currency,
		&payment.Status,
		&createdAt,
	); err != nil {
		return nil, err
	}

	payment.Currency = currency.String
	if payment.Currency == "" {
		payment.Currency = "IDR"
	}
	payment.CreatedAt = createdAt.Format(time.RFC3339)

	return payment, nil
}

func (u *PaymentRepositoryImpl) GetByID(ctx context.Context, paymentID int, db *sql.DB) (*proto.PaymentResponse, error) {
	SQL := `SELECT id, order_id, amount, currency, payment_method, payment_channel, 
			gateway_name, gateway_transaction_id, gateway_order_id, gateway_token, 
			gateway_redirect_url, va_number, qr_code_url, status, created_at, paid_at, expired_at 
			FROM payments WHERE id = $1`

	row := db.QueryRowContext(ctx, SQL, paymentID)

	payment := &proto.PaymentResponse{}
	var (
		currency             sql.NullString
		paymentMethod        sql.NullString
		paymentChannel       sql.NullString
		gatewayName          sql.NullString
		gatewayTransactionID sql.NullString
		gatewayOrderID       sql.NullString
		gatewayToken         sql.NullString
		gatewayRedirectURL   sql.NullString
		vaNumber             sql.NullString
		qrCodeURL            sql.NullString
		createdAt            time.Time
		paidAt               sql.NullTime
		expiredAt            sql.NullTime
	)

	if err := row.Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Amount,
		&currency,
		&paymentMethod,
		&paymentChannel,
		&gatewayName,
		&gatewayTransactionID,
		&gatewayOrderID,
		&gatewayToken,
		&gatewayRedirectURL,
		&vaNumber,
		&qrCodeURL,
		&payment.Status,
		&createdAt,
		&paidAt,
		&expiredAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	payment.Currency = currency.String
	payment.PaymentMethod = paymentMethod.String
	payment.PaymentChannel = paymentChannel.String
	payment.GatewayName = gatewayName.String
	payment.GatewayTransactionId = gatewayTransactionID.String
	payment.GatewayOrderId = gatewayOrderID.String
	payment.GatewayToken = gatewayToken.String
	payment.GatewayRedirectUrl = gatewayRedirectURL.String
	payment.VaNumber = vaNumber.String
	payment.QrCodeUrl = qrCodeURL.String
	payment.CreatedAt = createdAt.Format(time.RFC3339)

	if paidAt.Valid {
		payment.PaidAt = paidAt.Time.Format(time.RFC3339)
	}
	if expiredAt.Valid {
		payment.ExpiredAt = expiredAt.Time.Format(time.RFC3339)
	}

	return payment, nil
}

func (u *PaymentRepositoryImpl) GetByOrderID(ctx context.Context, orderID int, db *sql.DB) (*proto.PaymentResponse, error) {
	SQL := `SELECT id, order_id, amount, currency, payment_method, payment_channel, 
			gateway_name, gateway_transaction_id, gateway_order_id, gateway_token, 
			gateway_redirect_url, va_number, qr_code_url, status, created_at, paid_at, expired_at 
			FROM payments WHERE order_id = $1`

	row := db.QueryRowContext(ctx, SQL, orderID)

	payment := &proto.PaymentResponse{}
	var (
		currency             sql.NullString
		paymentMethod        sql.NullString
		paymentChannel       sql.NullString
		gatewayName          sql.NullString
		gatewayTransactionID sql.NullString
		gatewayOrderID       sql.NullString
		gatewayToken         sql.NullString
		gatewayRedirectURL   sql.NullString
		vaNumber             sql.NullString
		qrCodeURL            sql.NullString
		createdAt            time.Time
		paidAt               sql.NullTime
		expiredAt            sql.NullTime
	)

	if err := row.Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Amount,
		&currency,
		&paymentMethod,
		&paymentChannel,
		&gatewayName,
		&gatewayTransactionID,
		&gatewayOrderID,
		&gatewayToken,
		&gatewayRedirectURL,
		&vaNumber,
		&qrCodeURL,
		&payment.Status,
		&createdAt,
		&paidAt,
		&expiredAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found for this order")
		}
		return nil, err
	}

	payment.Currency = currency.String
	payment.PaymentMethod = paymentMethod.String
	payment.PaymentChannel = paymentChannel.String
	payment.GatewayName = gatewayName.String
	payment.GatewayTransactionId = gatewayTransactionID.String
	payment.GatewayOrderId = gatewayOrderID.String
	payment.GatewayToken = gatewayToken.String
	payment.GatewayRedirectUrl = gatewayRedirectURL.String
	payment.VaNumber = vaNumber.String
	payment.QrCodeUrl = qrCodeURL.String
	payment.CreatedAt = createdAt.Format(time.RFC3339)

	if paidAt.Valid {
		payment.PaidAt = paidAt.Time.Format(time.RFC3339)
	}
	if expiredAt.Valid {
		payment.ExpiredAt = expiredAt.Time.Format(time.RFC3339)
	}

	return payment, nil
}

func (u *PaymentRepositoryImpl) GetByGatewayOrderID(ctx context.Context, gatewayOrderID string, db *sql.DB) (*proto.PaymentResponse, error) {
	SQL := `SELECT id, order_id, amount, currency, payment_method, payment_channel, 
			gateway_name, gateway_transaction_id, gateway_order_id, gateway_token, 
			gateway_redirect_url, va_number, qr_code_url, status, created_at, paid_at, expired_at 
			FROM payments WHERE gateway_order_id = $1`

	row := db.QueryRowContext(ctx, SQL, gatewayOrderID)

	payment := &proto.PaymentResponse{}
	var (
		currency             sql.NullString
		paymentMethod        sql.NullString
		paymentChannel       sql.NullString
		gatewayName          sql.NullString
		gatewayTransactionID sql.NullString
		gatewayOrderIDVal    sql.NullString
		gatewayToken         sql.NullString
		gatewayRedirectURL   sql.NullString
		vaNumber             sql.NullString
		qrCodeURL            sql.NullString
		createdAt            time.Time
		paidAt               sql.NullTime
		expiredAt            sql.NullTime
	)

	if err := row.Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Amount,
		&currency,
		&paymentMethod,
		&paymentChannel,
		&gatewayName,
		&gatewayTransactionID,
		&gatewayOrderIDVal,
		&gatewayToken,
		&gatewayRedirectURL,
		&vaNumber,
		&qrCodeURL,
		&payment.Status,
		&createdAt,
		&paidAt,
		&expiredAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment not found by gateway order ID")
		}
		return nil, err
	}

	payment.Currency = currency.String
	payment.PaymentMethod = paymentMethod.String
	payment.PaymentChannel = paymentChannel.String
	payment.GatewayName = gatewayName.String
	payment.GatewayTransactionId = gatewayTransactionID.String
	payment.GatewayOrderId = gatewayOrderIDVal.String
	payment.GatewayToken = gatewayToken.String
	payment.GatewayRedirectUrl = gatewayRedirectURL.String
	payment.VaNumber = vaNumber.String
	payment.QrCodeUrl = qrCodeURL.String
	payment.CreatedAt = createdAt.Format(time.RFC3339)

	if paidAt.Valid {
		payment.PaidAt = paidAt.Time.Format(time.RFC3339)
	}
	if expiredAt.Valid {
		payment.ExpiredAt = expiredAt.Time.Format(time.RFC3339)
	}

	return payment, nil
}

func (u *PaymentRepositoryImpl) UpdatePaymentGateway(ctx context.Context, payment *proto.PaymentResponse, db *sql.DB) error {
	SQL := `UPDATE payments SET 
			payment_method = $1, 
			payment_channel = $2,
			gateway_order_id = $3,
			gateway_token = $4,
			gateway_redirect_url = $5,
			va_number = $6,
			qr_code_url = $7,
			expired_at = $8,
			status = $9
			WHERE order_id = $10`

	var expiredAt interface{}
	if payment.ExpiredAt != "" {
		t, err := time.Parse(time.RFC3339, payment.ExpiredAt)
		if err == nil {
			expiredAt = t
		}
	}

	_, err := db.ExecContext(ctx, SQL,
		payment.PaymentMethod,
		payment.PaymentChannel,
		payment.GatewayOrderId,
		payment.GatewayToken,
		payment.GatewayRedirectUrl,
		payment.VaNumber,
		payment.QrCodeUrl,
		expiredAt,
		payment.Status,
		payment.OrderId,
	)

	return err
}

func (u *PaymentRepositoryImpl) UpdatePaymentStatus(ctx context.Context, orderID int, status string, transactionID string, db *sql.DB) error {
	loc := time.FixedZone("WIB", 7*60*60)
	now := time.Now().In(loc)

	var SQL string
	var args []interface{}

	if status == "success" || status == "paid" {
		SQL = `UPDATE payments SET status = $1, gateway_transaction_id = $2, paid_at = $3 WHERE order_id = $4`
		args = []interface{}{status, transactionID, now, orderID}
	} else {
		SQL = `UPDATE payments SET status = $1, gateway_transaction_id = $2 WHERE order_id = $3`
		args = []interface{}{status, transactionID, orderID}
	}

	_, err := db.ExecContext(ctx, SQL, args...)
	return err
}
