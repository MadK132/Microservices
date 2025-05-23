package dto

import (
	"net/http"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"github.com/gin-gonic/gin"
)

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

type ProductUpdateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	Category    *string  `json:"category"`
}

type ProductCreateResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at"`
}

type ProductUpdateResponse = ProductCreateResponse

func FromProductCreateRequest(ctx *gin.Context) (model.Product, error) {
	var req ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return model.Product{}, err
	}

	now := time.Now()

	return model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsDeleted:   false,
	}, nil
}

func FromProductUpdateRequest(ctx *gin.Context) (model.ProductUpdate, error) {
	var req ProductUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.ProductUpdate{}, err
	}

	update := model.ProductUpdate{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	return update, nil
}

func ToProductCreateResponse(p model.Product) ProductCreateResponse {
	return ProductCreateResponse{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}
