package server

import "github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/grpc/server/frontend"

type ProductUsecase interface {
	frontend.ProductUsecase
}

type DiscountUsecase interface {
	frontend.DiscountUsecase
}

type ReviewUsecase interface {
	frontend.ReviewUsecase
}