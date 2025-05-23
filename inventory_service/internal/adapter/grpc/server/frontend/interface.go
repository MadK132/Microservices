package frontend

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductUsecase interface {
	Create(ctx context.Context, request model.Product) (model.Product, error)
	Update(ctx context.Context, request model.ProductUpdate) (model.Product, error)
	GetByID(ctx context.Context, filter model.ProductFilter) (model.Product, error)
	Delete(ctx context.Context, filter model.ProductFilter) error
	GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error)
	GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Product, error) 
}
type DiscountUsecase interface {
	CreateDiscnout(ctx context.Context, discount model.Discount) (model.Discount, error)
	GetAllProductsWithDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Product, error)
	GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error)
	DeleteDiscount(ctx context.Context, filter model.DiscountFilter) error
}

type ReviewUsecase interface{
	CreateReview(ctx context.Context, review model.Review) (model.Review, error)
	GetReviewByID(ctx context.Context, filter model.ReviewFilter) (model.Review, error)
	UpdateReview(ctx context.Context, request model.ReviewUpdate) (model.Review, error)
}