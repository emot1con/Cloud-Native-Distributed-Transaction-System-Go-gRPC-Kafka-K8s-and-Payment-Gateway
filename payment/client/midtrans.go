package client

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

var SnapClient snap.Client
var ServerKey string

func InitMidtransClient() {
	ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	if ServerKey == "" {
		logrus.Warn("MIDTRANS_SERVER_KEY not set, using default sandbox key")
		ServerKey = "SB-Mid-server-YOUR_SERVER_KEY"
	}

	SnapClient.New(ServerKey, midtrans.Sandbox)

	logrus.Infof("Midtrans initialized with client key: %s", clientKey)
}

// CreateSnapTransaction creates a Snap transaction and returns the token
func CreateSnapTransaction(orderID string, amount int64, customerName, customerEmail, customerPhone string) (*snap.Response, error) {
	// Create Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
			Phone: customerPhone,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "hour",
			Duration:  24,
		},
	}

	logrus.Infof("Creating Snap transaction for order: %s, amount: %d", orderID, amount)

	snapResp, err := SnapClient.CreateTransaction(req)
	if err != nil {
		logrus.Errorf("Failed to create Snap transaction: %v", err)
		return nil, err
	}

	logrus.Infof("Snap transaction created: token=%s, redirect_url=%s", snapResp.Token, snapResp.RedirectURL)
	return snapResp, nil
}

// VerifySignature verifies the webhook signature from Midtrans
func VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	// Signature = SHA512(order_id + status_code + gross_amount + ServerKey)
	data := orderID + statusCode + grossAmount + ServerKey
	hash := sha512.Sum512([]byte(data))
	calculatedSignature := hex.EncodeToString(hash[:])

	isValid := calculatedSignature == signatureKey
	if !isValid {
		logrus.Warnf("Invalid signature for order %s: expected %s, got %s", orderID, calculatedSignature, signatureKey)
	}

	return isValid
}

// MapTransactionStatus maps Midtrans status to our internal status
func MapTransactionStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return "paid"
		}
		return "pending"
	case "settlement":
		return "paid"
	case "pending":
		return "pending"
	case "deny", "cancel", "expire":
		return "failed"
	case "refund", "partial_refund":
		return "refunded"
	default:
		return "pending"
	}
}

// GenerateOrderID generates a unique order ID for Midtrans
func GenerateOrderID(paymentID int32) string {
	return fmt.Sprintf("PAY-%d-%d", paymentID, time.Now().Unix())
}
