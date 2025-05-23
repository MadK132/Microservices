package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/pkg/redis"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const keyPrefix = "product:%s"

type Product struct {
	product *redis.Product
	ttl     time.Duration
}

func NewProduct(product *redis.Product, ttl time.Duration) *Product {
	return &Product{
		product: product,
		ttl:     ttl,
	}
}

func (p *Product) Set(ctx context.Context, product model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}
	fmt.Println("Set in Redis:", product.ID.Hex())
	return p.product.Unwrap().Set(ctx, p.key(product.ID), data, p.ttl).Err()
}

func (p *Product) SetMany(ctx context.Context, products []model.Product) error {
	pipe := p.product.Unwrap().Pipeline()
	for _, product := range products {
		data, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product: %w", err)
		}
		pipe.Set(ctx, p.key(product.ID), data, p.ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set many products: %w", err)
	}
	return nil
}

func (p *Product) Get(ctx context.Context, id string) (model.Product, error) {
	fmt.Println("Try get from Redis:", id)
	data, err := p.product.Unwrap().Get(ctx, fmt.Sprintf(keyPrefix, id)).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return model.Product{}, nil
		}
		return model.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	var product model.Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return product, nil
}

func (p *Product) GetAll(ctx context.Context) ([]model.Product, error) {
	keys, err := p.product.Unwrap().Keys(ctx, "product:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	var products []model.Product
	for _, key := range keys {
		data, err := p.product.Unwrap().Get(ctx, key).Bytes()
		if err != nil {
			if err == goredis.Nil {
				continue
			}
			return nil, fmt.Errorf("failed to get product: %w", err)
		}

		var product model.Product
		err = json.Unmarshal(data, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *Product) key(id primitive.ObjectID) string {
	return fmt.Sprintf(keyPrefix, id.Hex())
}
