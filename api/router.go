package api

import (
	"errors"
	"fmt"
	"net/http"
	"rent-car/api/handler"
	"rent-car/pkg/logger"
	"rent-car/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "rent-car/api/docs"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceManager, log logger.ILogger) *gin.Engine {
	h := handler.NewStrg(services, log)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/customer/login", h.LoginCustomer)
	r.POST("/customer/register", h.CustomerRegister)
	r.POST("/customer/register-confirm", h.CustomerRegisterConfirm)
	r.POST("/customer", h.CreateCustomer)

	//r.Use(authMiddleware)
	//r.Use(logMiddleware)

	r.POST("/car", h.CreateCar)
	r.PUT("/car/:id", h.UpdateCar)
	r.GET("/car/:id", h.GetCarByID)
	r.GET("/car", h.GetAllCars)
	r.GET("car/available", h.GetAvailableCars)
	r.DELETE("/car/:id", h.DeleteCar)

	r.PUT("/customer/:id", h.UpdateCustomer)
	r.PATCH("/customer", h.ChangePasswordCustomer)
	r.GET("/customer/:id", h.GetCustomerByID)
	r.GET("/customer", h.GetAllCustomers)
	r.GET("/customer/cars", h.GetCustomerCars)
	r.DELETE("/customer/:id", h.DeleteCustomer)

	r.POST("/order", h.CreateOrder)
	r.PUT("/order/:id", h.UpdateOrder)
	r.PATCH("/order", h.UpdateOrderStatus)
	r.GET("/order/:id", h.GetOrderByID)
	r.GET("/order", h.GetAllOrders)
	r.DELETE("/order/:id", h.DeleteOrder)

	return r
}

func authMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	c.Next()
}

func logMiddleware(c *gin.Context) {
	headers := c.Request.Header

	for key, values := range headers {
		for _, v := range values {
			fmt.Printf("Header: %v, Value: %v\n", key, v)
		}
	}

	c.Next()
}
