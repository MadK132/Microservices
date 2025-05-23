package dto

import (
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FromProtoReviewCreate(req *inventorypb.CreateReviewRequest) model.Review{
	var productIDs primitive.ObjectID
	for _, id := range req.GetProductId() {
		objID, err := primitive.ObjectIDFromHex(string(id))
		if err != nil {
			return model.Review{}
		}
		productIDs = objID
	}

	var UserID primitive.ObjectID
	for _, id := range req.GetUserId() {
		objID, err := primitive.ObjectIDFromHex(string(id))
		if err != nil {
			return model.Review{}
		}
		UserID = objID
	}
	return model.Review{
		ProductID: productIDs,
		UserID: UserID,
		Rating: float64(req.GetRating()),
		Comment: req.GetComment(),
	}
}

func ToProtoReview(d model.Review) *inventorypb.Review{
	return &inventorypb.Review{
		Id: d.ID.Hex(),
		ProductId: d.ProductID.Hex(),
		UserId: d.UserID.Hex(),
		Rating: float32(d.Rating),
		Comment: d.Comment,
	}
}