package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/nats/producer"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	repo ProductRepo
	producer *producer.InventoryProducer
	inMemoryCache ProductCache
	redisCache RedisCache
}

func NewProduct(repo ProductRepo, producer *producer.InventoryProducer, inMemoryCache ProductCache, redisCache RedisCache) *Product {
	return &Product{
		repo: repo,
		producer: producer,
		inMemoryCache: inMemoryCache,
		redisCache: redisCache,
	}
}

func (p *Product) Create(ctx context.Context, request model.Product) (model.Product, error) {
	request.ID = primitive.NewObjectID()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.IsDeleted = false
	if err := p.repo.Create(ctx, request); err != nil {
		return model.Product{}, err
		
	}

	p.inMemoryCache.Set(request)
	err := p.redisCache.Set(ctx, request)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to set product in redis: %w", err)
	}
	_ = p.producer.Push(ctx, request.ID.Hex(), "created")

	return request, nil
}

func (p *Product) GetByID(ctx context.Context, request model.ProductFilter) (model.Product, error) {
	if request.ID == nil {
		return model.Product{}, model.ErrProductNotFound
	}

	idStr := request.ID.Hex()

	if product, ok := p.inMemoryCache.Get(idStr); ok {
		return product, nil
	}

	product, err := p.redisCache.Get(ctx, idStr)
	if err == nil && product.ID != primitive.NilObjectID {
		p.inMemoryCache.Set(product)
		return product, nil
	}

	product, err = p.repo.GetByID(ctx, request)
	if err != nil {
		return model.Product{}, err
	}

	p.inMemoryCache.Set(product)
	_ = p.redisCache.Set(ctx, product)

	return product, nil
}

func (p *Product) Update(ctx context.Context, request model.ProductUpdate) (model.Product, error) {
	if request.ID == nil {
		return model.Product{}, model.ErrProductNotFound
	}

	if err := p.repo.Update(ctx, model.ProductFilter{ID: request.ID}, request); err != nil {
		return model.Product{}, err
	}

	updated, err := p.repo.GetByID(ctx, model.ProductFilter{ID: request.ID})
	if err != nil {
		return model.Product{}, err
	}

	_ = p.producer.Push(ctx, request.ID.Hex(), "updated")

	return updated, nil
}

func (p *Product) Delete(ctx context.Context, request model.ProductFilter) error {
	if request.ID == nil {
		return model.ErrProductNotFound
	}
	_ = p.producer.Push(ctx, request.ID.Hex(), "deleted")
	return p.repo.Delete(ctx, request)
}

func (p *Product) GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	if filter.ID == nil && filter.Name == nil {
		return p.inMemoryCache.GetAll(), nil
	}
	return p.repo.GetAll(ctx, filter)
}

func (p *Product) GetByIDs(ctx context.Context, ids[] primitive.ObjectID) ([]model.Product, error) {
	return p.repo.GetByIDs(ctx, ids)
}

func (p *Product) StartCacheAutoRefresh(ctx context.Context) {
	go func() {
		for {
			time.Sleep(12 * time.Hour)

			products, err := p.repo.GetAll(ctx, model.ProductFilter{})
			if err != nil {
				fmt.Println("Failed to refresh cache:", err)
				continue
			}
			p.inMemoryCache.SetMany(products)
			fmt.Println("In-memory product cache refreshed.")
		}
	}()
}
