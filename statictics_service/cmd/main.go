package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/statistics_service/internal/adapter/grpc/server"
	"github.com/recktt77/statistics_service/internal/adapter/nats/handler"
	_"github.com/recktt77/statistics_service/internal/repository"
	"github.com/recktt77/statistics_service/internal/usecase"
	mongocon "github.com/recktt77/statistics_service/internal/adapter/mongo"
	mongorepo "github.com/recktt77/statistics_service/internal/adapter/mongo"
)

func main() {
	nc, _ := nats.Connect("nats://localhost:4222")
	defer nc.Close()

	db, err := mongocon.Connect("mongodb://localhost:27017", "statistics_db")
	if err != nil {
		log.Fatal("Mongo error:", err)
	}

	mongoRepo := mongorepo.NewMongoStatsRepo(db)
	statsUC := usecase.NewStatsUsecase(mongoRepo)


	consumer := handler.NewNATSConsumer(statsUC, statsUC)
	consumer.Subscribe(nc)
	go func() {
		grpcSrv := server.NewGRPCServer(mongoRepo)
		_ = grpcSrv.Run(":8083")
	}()

	select {}
}