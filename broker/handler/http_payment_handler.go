package handler

import (
	"broker/auth"
	"broker/proto"
	"broker/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	repo repository.PaymentRepository
}

func NewPaymentHandler(repo repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		repo: repo,
	}
}

func (u *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	paymentRoutes := r.Group("/payment")
	paymentRoutes.Use(auth.ProtectedEndpoint())

	paymentRoutes.POST("/transaction", u.Transaction)
	paymentRoutes.GET("/order/:order_id", u.GetPaymentByOrderId)
}

func (u *PaymentHandler) Transaction(c *gin.Context) {
	var payload proto.PaymentTransaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := u.repo.Transaction(&payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "transaction success, your order is being processed"})
}

func (u *PaymentHandler) GetPaymentByOrderId(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := u.repo.GetPaymentByOrderId(int32(orderID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, payment)
}
