package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/app"
)

func main() {
	err := godotenv.Load("C:\\Users\\1\\OneDrive\\Рабочий стол\\University\\Micriservices\\inventory_service\\.env")
	if err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}

	log.Println("MONGO_DB_URI:", os.Getenv("MONGO_DB_URI"))

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	natsConn, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("Не удалось подключиться к NATS: %v", err)
	}
	defer natsConn.Close()

	ctx := context.Background()
	application, err := app.New(ctx, cfg, natsConn)
	if err != nil {
		log.Fatalf("Ошибка при создании приложения: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка выполнения приложения: %v", err)
	}
}
