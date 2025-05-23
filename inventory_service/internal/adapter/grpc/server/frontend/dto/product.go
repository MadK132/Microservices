package dto

import (
	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

func FromProtoCreate(req *inventorypb.CreateProductRequest) model.Product {
	return model.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       float64(req.GetPrice()),
		Stock:       int(req.GetStock()),
		Category:    req.GetCategory(),
	}
}

func ToProto(p model.Product) *inventorypb.Product {
	return &inventorypb.Product{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       float64(p.Price),
		Stock:       int32(p.Stock),
		Category:    p.Category,
	}
}
