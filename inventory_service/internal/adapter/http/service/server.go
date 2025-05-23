package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/http/service/handler"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string

	productHandler *handler.Product
	discountHandler *handler.Discount
}

func (a *API) setupRoutes() {
	v1 := a.server.Group("/api/v1")
	{
		product := v1.Group("/products")
		{
			product.POST("/", a.productHandler.Create)
			product.GET("/:id", a.productHandler.GetByID)
			product.PATCH("/:id", a.productHandler.Update)
			product.DELETE("/:id", a.productHandler.Delete)
			product.GET("/", a.productHandler.GetAll)
		}
		discount := v1.Group("/discount")
		{
			discount.POST("/", a.discountHandler.CreateDiscnout)
			discount.GET("/", a.discountHandler.GetAllExistingDiscounts)
			discount.GET("/products", a.discountHandler.GetAllProductsWithDiscounts)
			discount.DELETE("/:id", a.discountHandler.DeleteDiscount)
		}
	}

	// Swagger documentation
	a.server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func New(cfg config.Server, productUsecase ProductUsecase, discountUsecase DiscountUsecase) *API {
	//setting the Gin mode
	gin.SetMode(cfg.HTTPServer.Mode)
	// creating a new gin engine
	server := gin.New()

	//applying the middleware
	server.Use(gin.Recovery())

	// building proucts
	productHandler := handler.NewProduct(productUsecase)
	discountHandler := handler.NewDiscount(discountUsecase)

	api := &API{
		server:         server,
		cfg:            cfg.HTTPServer,
		addr:           fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		productHandler: productHandler,
		discountHandler: discountHandler,
	}

	api.setupRoutes()
	return api
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server starting on: %v", a.addr)
		if err := a.server.Run(a.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) Stop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Shutting signal received: %v", sig.String())

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")

	log.Println("HTTP server stopped successfully")

	return nil
}
