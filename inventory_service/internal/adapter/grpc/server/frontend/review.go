package frontend

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	svc "github.com/recktt77/proto-definitions/gen/inventory"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct{
	svc.UnimplementedReviewServiceServer

	uc ReviewUsecase
}

func NewReview(uc ReviewUsecase) *Review{
	return &Review{
		uc:uc,
	}
}

func (r *Review) CreateReview(ctx context.Context, req *svc.CreateReviewRequest) (*svc.Review, error){
	var productIDs primitive.ObjectID
	id := req.GetProductId() 
	objID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}
	productIDs = objID

	var UserID primitive.ObjectID
	uid := req.GetUserId() 
	uobjID, err := primitive.ObjectIDFromHex(string(uid))
	if err != nil {
		return nil, err
	}
	UserID = uobjID
	domainProduct := model.Review{
		ProductID: productIDs,
		UserID: UserID,
		Rating: float64(req.GetRating()),
		Comment: req.GetComment(),
	}

	created, err := r.uc.CreateReview(ctx, domainProduct)
	if err != nil {
		return nil, err
	}

	return &svc.Review{
		Id: created.ID.Hex(),
		ProductId: created.ProductID.Hex(),
		UserId: created.UserID.Hex(),
		Rating: float32(created.Rating),
		Comment: created.Comment,
	},nil 
}

func (r *Review) GetReviewByID(ctx context.Context, req *svc.ReviewID) (*svc.Review, error){
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	review, err := r.uc.GetReviewByID(ctx, model.ReviewFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	return &svc.Review{
		Id: review.ID.Hex(),
		ProductId: review.ProductID.Hex(),
		UserId: review.UserID.Hex(),
		Rating: float32(review.Rating),
		Comment: review.Comment,
	},nil 
}

func (r *Review) UpdateReview(ctx context.Context, req *svc.UpdateReviewRequest) (*svc.Review, error){
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	var productIDs primitive.ObjectID
	productID := req.GetProductId()
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}
	productIDs = objID

	var UserID primitive.ObjectID
	uid := *req.UserId 
	uobjID, err := primitive.ObjectIDFromHex(string(uid))
		if err != nil {
			return nil, err
		}
		UserID = uobjID

	update := model.ReviewUpdate{
		ID: &id,
		ProductID: &productIDs,
		UserID: &UserID,
		Rating: func(r *float32) *float64 {
			if r == nil {
				return nil
			}
			val := float64(*r)
			return &val
		}(req.Rating),
		Comment: req.Comment,
	}
	updated, err := r.uc.UpdateReview(ctx, update)
	if err != nil {
		return nil, err
	}

	return &svc.Review{
		Id: updated.ID.Hex(),
		ProductId: updated.ProductID.Hex(),
		UserId: updated.UserID.Hex(),
		Rating: float32(updated.Rating),
		Comment: func(c *string) string {
			if c == nil {
				return ""
			}
			return *c
		}(update.Comment),
	}, nil 
}