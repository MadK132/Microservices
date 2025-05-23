package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/grpc/server"
	inmemorycache "github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/inmemory"
	mongorepo "github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/nats/producer"
	rediscache "github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/redis"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/usecase"
	mongocon "github.com/recktt77/Microservices-First-/inventory_service/pkg/mongo"
	redispkg "github.com/recktt77/Microservices-First-/inventory_service/pkg/redis"
)

const (
	serviceName     = "product-service"
	shutdownTimeout = 30 * time.Second
)

type App struct {
	grpcServer *server.API
	redisClient      *redispkg.Product
	productCache     *rediscache.Product
	inMemoryCache    *inmemorycache.Product
}

func New(ctx context.Context, cfg *config.Config, natsConn *nats.Conn) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	mongoDB, err := mongocon.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	redisClient, err := redispkg.NewProduct(ctx, redispkg.Config(cfg.Redis))
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	productDAO := dao.NewProductDAO(mongoDB.Conn)
	discountDAO := dao.NewDiscountDAO(mongoDB.Conn)
	reviewDAO := dao.NewReviewDAO(mongoDB.Conn)

	productRepo := mongorepo.NewProductRepository(productDAO)
	discountRepo := mongorepo.NewDiscountRepository(discountDAO)
	reviewRepo := mongorepo.NewReviewRepository(reviewDAO)

	productMemoryCache := inmemorycache.NewProduct()
	productRedisCache := rediscache.NewProduct(redisClient, cfg.Cache.ProductTTL)

	products, err := productRepo.GetAll(ctx, model.ProductFilter{})
	if err == nil {
		productMemoryCache.SetMany(products)
		log.Println("in-memory product cache initialized from DB")
	} else {
		log.Println("failed to init in-memory cache:", err)
	}

	prod := producer.NewInventoryProducer(natsConn)
	productUC := usecase.NewProduct(productRepo, prod, productMemoryCache, productRedisCache)
	discountUC := usecase.NewDiscount(discountRepo, productRepo)
	reviewUC := usecase.NewReview(reviewRepo)

	productUC.StartCacheAutoRefresh(ctx)

	grpcSrv := server.New(cfg.Server.GRPCServer, productUC, discountUC, reviewUC)

	return &App{
	grpcServer:     grpcSrv,
	redisClient:    redisClient,
	productCache:   productRedisCache,
	inMemoryCache:  productMemoryCache,
	}, nil

}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	go a.grpcServer.Run(ctx, errCh)
	log.Println(fmt.Sprintf("service %v started (gRPC)", serviceName))

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("gRPC server error: %w", err)
	case sig := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Shutting down...", sig))
		return a.Close()
	}
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.redisClient.Close(); err != nil {
		log.Println("failed to close redis:", err)
	}

	if err := a.grpcServer.Stop(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out after %v", shutdownTimeout)
	}

	log.Println("graceful shutdown completed successfully")
	return nil
}
