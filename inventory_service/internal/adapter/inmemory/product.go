package inmemory

import (
	"fmt"
	"sync"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

type Product struct {
	mu      sync.RWMutex
	storage map[string]model.Product
}

func NewProduct() *Product {
	return &Product{
		storage: make(map[string]model.Product),
	}
}

func (c *Product) Set(product model.Product) {
	fmt.Println("Set in memory:", product.ID.Hex())
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[product.ID.Hex()] = product
}

func (c *Product) SetMany(products []model.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range products {
		c.storage[p.ID.Hex()] = p
	}
}

func (c *Product) Get(id string) (model.Product, bool) {
	fmt.Println("Try get from memory:", id)
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.storage[id]
	return val, ok
}

func (c *Product) GetAll() []model.Product {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []model.Product
	for _, val := range c.storage {
		result = append(result, val)
	}
	return result
}
