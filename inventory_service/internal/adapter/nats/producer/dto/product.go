package dto

import (
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	events "github.com/recktt77/proto-definitions/gen/inventory"
	
)

func FromProduct(product model.Product) events.Product{
	return events.Product{
		Id: product.ID.Hex(),
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
		Stock: int32(product.Stock),
		Category: product.Category,
		IsDeleted: product.IsDeleted,
	}
}