package dto

import (

	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FromProtoCreateDiscountRequest(req *inventorypb.CreateDiscountRequest) (model.Discount, error) {
	var productIDs []primitive.ObjectID
	for _, id := range req.GetApplicableProductIds() {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return model.Discount{}, err
		}
		productIDs = append(productIDs, objID)
	}


	return model.Discount{
		Name:               req.GetName(),
		Description:        req.GetDescription(),
		DiscountPercentage: float64(req.GetDiscountPercentage()),
		ApplicableProducts: productIDs,
		IsActive:           true,
	}, nil
}

func ToProtoDiscount(d model.Discount) *inventorypb.Discount {
	var ids []string
	for _, id := range d.ApplicableProducts {
		ids = append(ids, id.Hex())
	}

	return &inventorypb.Discount{
		Id:                  d.ID.Hex(),
		Name:                d.Name,
		Description:         d.Description,
		DiscountPercentage:  float64(d.DiscountPercentage),
		ApplicableProductIds: ids,
		IsActive:            d.IsActive,
	}
}
