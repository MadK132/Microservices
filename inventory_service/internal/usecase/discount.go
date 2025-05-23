package usecase

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	repo        DiscountRepo
	productRepo ProductRepo 
}

func NewDiscount(repo DiscountRepo, productRepo ProductRepo) *Discount {
	return &Discount{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (d *Discount) CreateDiscnout(ctx context.Context, request model.Discount) (model.Discount, error){
	if err := d.repo.CreateDiscnout(ctx, request); err != nil{
		return model.Discount{}, err
	}
	return request, nil
}

func (d *Discount) GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error){
	discounts, err := d.repo.GetAllExistingDiscounts(ctx, filter)
	if err != nil{
		return []model.Discount{}, err
	}

	return discounts, err
}

func (d *Discount) GetAllProductsWithDiscounts(ctx context.Context, request model.DiscountFilter) ([]model.Product, error){
	discounts, err := d.repo.GetAllExistingDiscounts(ctx, request)
	if err != nil {
		return nil, err
	}

	productIDMap := make(map[primitive.ObjectID]struct{})
	for _,discount := range discounts {
		for _,pid := range discount.ApplicableProducts{
			productIDMap[pid]= struct{}{}
		}
	}

	var ids []primitive.ObjectID
	for pid := range productIDMap {
		ids = append(ids, pid)
	}

	products, err := d.productRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return products, nil

}

func (d *Discount) DeleteDiscount(ctx context.Context, request model.DiscountFilter) error{
	if request.ID == nil{
		return model.ErrDiscnoutNotFound
	}
	return d.repo.DeleteDiscount(ctx, request)
}