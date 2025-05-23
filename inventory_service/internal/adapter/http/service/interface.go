package service

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductUsecase interface {
	Create(ctx context.Context, product model.Product) (model.Product, error)
	GetByID(ctx context.Context, filter model.ProductFilter) (model.Product, error)
	Update(ctx context.Context, product model.ProductUpdate) (model.Product, error)
	Delete(ctx context.Context, filter model.ProductFilter) error
	GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error)
	GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Product, error)
}

type DiscountUsecase interface {
	CreateDiscnout(ctx context.Context, discount model.Discount) (model.Discount, error)
	GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error)
	GetAllProductsWithDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Product, error)
	DeleteDiscount(ctx context.Context, filter model.DiscountFilter) error
}
