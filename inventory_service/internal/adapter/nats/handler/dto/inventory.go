package dto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToProduct(msg *nats.Msg) (model.Product, error){
	var pbProduct inventorypb.Product
	err := proto.Unmarshal(msg.Data, &pbProduct)
	if err != nil {
		return model.Product{}, fmt.Errorf("proto.Unmarshall: %w", err)
	}

	return model.Product{
		ID: func() primitive.ObjectID {
			id, err := primitive.ObjectIDFromHex(pbProduct.Id)
			if err != nil {
				return primitive.NilObjectID
			}
			return id
		}(),
		Name: pbProduct.Name,
		Description: pbProduct.Description,
		Price: pbProduct.Price,
		Stock: int(pbProduct.Stock),
		Category: pbProduct.Category,
		IsDeleted: pbProduct.IsDeleted,
	}, nil
}