package mongo

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

type reviewRepository struct {
	dao *dao.ReviewDAO
}

func NewReviewRepository(dao *dao.ReviewDAO) *reviewRepository {
	return &reviewRepository{dao: dao}
}

func (r *reviewRepository) CreateReview(ctx context.Context, review model.Review) error {
	_, err := r.dao.CreateReview(ctx, review)
	return err
}

func (r *reviewRepository) GetReviewByID(ctx context.Context, filter model.ReviewFilter) (model.Review, error) {
	if filter.ID == nil {
		return model.Review{}, model.ErrReviewNotFound
	}
	return r.dao.GetReviewByID(ctx, *filter.ID)
}

func (r *reviewRepository) UpdateReview(ctx context.Context, filter model.ReviewFilter, update model.ReviewUpdate) error {
	if filter.ID == nil {
		return model.ErrReviewNotFound
	}
	return r.dao.UpdateReview(ctx, *filter.ID, update)
} 