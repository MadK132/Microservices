package mongo

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productRepository struct {
	dao *dao.ProductDAO
}

func NewProductRepository(dao *dao.ProductDAO) *productRepository {
	return &productRepository{dao: dao}
}

func (r *productRepository) Create(ctx context.Context, product model.Product) error {
	_, err := r.dao.Create(ctx, product)
	return err
}

func (r *productRepository) GetByID(ctx context.Context, filter model.ProductFilter) (model.Product, error) {
	if filter.ID == nil {
		return model.Product{}, model.ErrProductNotFound
	}
	return r.dao.GetByID(ctx, *filter.ID)
}

func (r *productRepository) Update(ctx context.Context, filter model.ProductFilter, update model.ProductUpdate) error {
	if filter.ID == nil {
		return model.ErrProductNotFound
	}
	return r.dao.Update(ctx, *filter.ID, update)
}

func (r *productRepository) Delete(ctx context.Context, filter model.ProductFilter) error {
	if filter.ID == nil {
		return model.ErrProductNotFound
	}
	return r.dao.Delete(ctx, *filter.ID)
}

func (r *productRepository) GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	return r.dao.GetAll(ctx, filter)
}

func (r *productRepository) GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Product, error){
	if len(ids) == 0{
		return []model.Product{}, model.ErrProductNotFound
	}
	return r.dao.GetByIDs(ctx, ids)
}