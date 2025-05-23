package mongo

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type discountRepository struct {
	dao *dao.DiscountDAO
}


func NewDiscountRepository(dao *dao.DiscountDAO) *discountRepository {
	return &discountRepository{dao: dao}
}

func (r *discountRepository) CreateDiscnout(ctx context.Context, promotion model.Discount) error {
	_, err := r.dao.CreateDiscnout(ctx, promotion)
	return err
}

func (r *discountRepository) GetAllProductsWithDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Product, error) {
	products, err := r.dao.GetAllProductsWithDiscounts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *discountRepository) DeleteDiscount(ctx context.Context, filter model.DiscountFilter) error {
	if filter.ID == nil || *filter.ID == primitive.NilObjectID {
		return model.ErrDiscnoutNotFound
	}
	return r.dao.DeletePromotion(ctx, *filter.ID)
}

func (r *discountRepository) GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error) {
	return r.dao.GetAllExistingDiscounts(ctx, filter)
}

