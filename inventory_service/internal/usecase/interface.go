package usecase

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ProductRepo interface {
	Create(ctx context.Context, product model.Product) error
	GetByID(ctx context.Context, filter model.ProductFilter) (model.Product, error)
	Update(ctx context.Context, filter model.ProductFilter, update model.ProductUpdate) error
	Delete(ctx context.Context, filter model.ProductFilter) error
	GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error)
	GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Product, error)
}

type DiscountRepo interface {
	CreateDiscnout(ctx context.Context, discount model.Discount) error
	GetAllProductsWithDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Product, error)
	GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error)
	DeleteDiscount(ctx context.Context, filter model.DiscountFilter) error
}

type ReviewRepo interface {
	CreateReview(ctx context.Context, review model.Review) error
	GetReviewByID(ctx context.Context, filter model.ReviewFilter) (model.Review, error)
	UpdateReview(ctx context.Context, filter model.ReviewFilter, update model.ReviewUpdate) error
}
type ProductCache interface {
	Set(product model.Product)
	Get(id string) (model.Product, bool)
	GetAll() []model.Product
	SetMany(products []model.Product)
}

type RedisCache interface {
	Set(ctx context.Context, product model.Product) error
	Get(ctx context.Context, id string) (model.Product, error)
}
