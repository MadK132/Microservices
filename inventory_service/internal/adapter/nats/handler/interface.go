package handler

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

type ProductUsecase interface {
	Create(ctx context.Context, product model.Product) error
}

type DiscountUsecase interface {
	CreateDiscount(ctx context.Context, discount model.Discount) error
}

type ReviewUsecase interface{
	CreateReview(ctx context.Context, review model.Review) error
}
