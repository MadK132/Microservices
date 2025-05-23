package usecase

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

type Review struct {
	repo ReviewRepo
}

func NewReview(repo ReviewRepo) *Review {
	return &Review{
		repo: repo,
	}
}

func (r *Review) CreateReview(ctx context.Context, request model.Review) (model.Review, error){
	if err := r.repo.CreateReview(ctx, request); err != nil{
		return model.Review{}, err
	}

	return request, nil
}


func (r *Review) GetReviewByID(ctx context.Context, request model.ReviewFilter) (model.Review, error) {
	if request.ID == nil {
		return model.Review{}, model.ErrProductNotFound
	}

	review, err := r.repo.GetReviewByID(ctx, request)
	if err != nil {
		return model.Review{}, err
	}

	return review, nil
}
func (r *Review) UpdateReview(ctx context.Context, request model.ReviewUpdate) (model.Review, error){
	if request.ID == nil {
		return model.Review{}, model.ErrReviewNotFound
	}
	if err := r.repo.UpdateReview(ctx, model.ReviewFilter{ID: request.ID}, request); err != nil{
		return model.Review{}, err
	}
	updated, err := r.repo.GetReviewByID(ctx, model.ReviewFilter{ID: request.ID})
	if err != nil {
		return model.Review{}, err
	}

	return updated, nil 
}