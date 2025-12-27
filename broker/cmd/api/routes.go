package routes

import (
	"broker/handler"
	"broker/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	r := gin.Default()

	// CORS middleware - allow all origins for development
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	userRepo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(userRepo)
	userHandler.RegisterRoutes(r)

	productRepo := repository.NewProductRepository()
	productHandler := handler.NewProductHandler(productRepo)
	productHandler.RegisterRoutes(r)

	orderRepo := repository.NewOrderRepository()
	orderHandler := handler.NewOrderHandler(userRepo, productRepo, orderRepo)
	orderHandler.RegisterRoutes(r)

	paymentRepo := repository.NewPaymentRepository()
	paymentHandler := handler.NewPaymentHandler(paymentRepo)
	paymentHandler.RegisterRoutes(r)

	return r
}
