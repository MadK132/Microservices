// internal/adapter/nats/handler/consumer.go
package handler

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

type NATSConsumer struct {
	productUsecase ProductUsecase
	orderUsecase   OrderUsecase
}

func NewNATSConsumer(pu ProductUsecase, ou OrderUsecase) *NATSConsumer {
	return &NATSConsumer{
		productUsecase: pu,
		orderUsecase:   ou,
	}
}

func (c *NATSConsumer) Subscribe(nc *nats.Conn) {
	ctx := context.Background()

	_, err := nc.Subscribe("inventory.*", func(msg *nats.Msg) {
		log.Println("Получено событие из", msg.Subject)
		if err := c.productUsecase.ProcessInventoryEvent(ctx, msg); err != nil {
			log.Println("Ошибка обработки inventory события:", err)
		}
	})
	if err != nil {
		log.Fatal("Ошибка подписки на inventory.*:", err)
	}

	_, err = nc.Subscribe("orders.*", func(msg *nats.Msg) {
		log.Println("Получено событие из", msg.Subject)
		if err := c.orderUsecase.ProcessOrderEvent(ctx, msg); err != nil {
			log.Println("Ошибка обработки order события:", err)
		}
	})
	if err != nil {
		log.Fatal("Ошибка подписки на order.*:", err)
	}

	log.Println("Подписка на NATS активна (inventory.*, order.*)")
}
