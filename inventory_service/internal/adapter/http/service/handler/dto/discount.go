package dto

import (
	"net/http"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountCreateRequest struct {
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	DiscountPercentage float64              `json:"discount_percentage"`
	ApplicableProducts []primitive.ObjectID             `json:"applicable_products"`
	StartDate          *time.Time           `json:"start_date,omitempty"`
	EndDate            *time.Time           `json:"end_date,omitempty"`
}

type DiscountCreateResponse struct {
	ID                 primitive.ObjectID    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ApplicableProducts []primitive.ObjectID  `json:"applicable_products"`
	StartDate          string    `json:"startdate"`
	EndDate            string    `json:"enddate"`
	IsActive           bool      `json:"isactive"`
}

func FromDiscountCreateRequest(ctx *gin.Context) (model.Discount, error) {
	var req DiscountCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return model.Discount{}, err
	}

	now := time.Now()

	start := now
	if req.StartDate != nil {
		start = *req.StartDate
	}

	end := start.Add(7 * 24 * time.Hour)
	if req.EndDate != nil {
		end = *req.EndDate
	}

	return model.Discount{
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
		ApplicableProducts: req.ApplicableProducts,
		StartDate:          start,
		EndDate:            end,
		IsActive:           true,
	}, nil
}


func ToDiscountCreateResponse(p model.Discount) DiscountCreateResponse {
	return DiscountCreateResponse{
		ID:                 p.ID,
		Name:               p.Name,
		Description:        p.Description,
		DiscountPercentage: p.DiscountPercentage,
		ApplicableProducts: p.ApplicableProducts,
		StartDate:          p.StartDate.Format(time.RFC3339),
		EndDate:            p.EndDate.Format(time.RFC3339),
		IsActive:           p.IsActive,
	}
}

