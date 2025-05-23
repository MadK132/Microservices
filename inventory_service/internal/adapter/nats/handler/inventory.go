package handler

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/nats/handler/dto"
)

type Product struct {
	usecase ProductUsecase
}

func NewClient(usecase ProductUsecase) *Product{
	return &Product{usecase: usecase}
}

func (p *Product) Handler(ctx context.Context, msg *nats.Msg) error{
	product, err := dto.ToProduct(msg)
	if err != nil{
		log.Println("failed to convert Product NATS msg", err)

		return err
	}

	err = p.usecase.Create(ctx, product)

	if err != nil{
		log.Println("failed to create product", err)

		return err
	}

	return nil 

}