package frontend

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	svc "github.com/recktt77/proto-definitions/gen/inventory"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Discount struct {
	svc.UnimplementedDiscountServiceServer

	uc DiscountUsecase
}


func NewDiscount(uc DiscountUsecase) *Discount {
	return &Discount{
		uc: uc,
	}
}

func (d *Discount) CreateDiscount(ctx context.Context, req *svc.CreateDiscountRequest) (*svc.Discount, error) {
	var productIDs []primitive.ObjectID
	for _, id := range req.GetApplicableProductIds() {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		productIDs = append(productIDs, objID)
	}


	discount := model.Discount{
		Name:               req.GetName(),
		Description:        req.GetDescription(),
		DiscountPercentage: float64(req.GetDiscountPercentage()),
		ApplicableProducts: productIDs,
		IsActive:           true,
	}

	created, err := d.uc.CreateDiscnout(ctx, discount)
	if err != nil {
		return nil, err
	}

	var productStrIDs []string
	for _, id := range created.ApplicableProducts {
		productStrIDs = append(productStrIDs, id.Hex())
	}

	return &svc.Discount{
		Id:                   created.ID.Hex(),
		Name:                 created.Name,
		Description:          created.Description,
		DiscountPercentage:   float64(created.DiscountPercentage),
		ApplicableProductIds: productStrIDs,
		IsActive:             created.IsActive,
	}, nil
}

func (d *Discount) GetAllDiscounts(ctx context.Context, _ *emptypb.Empty) (*svc.DiscountList, error) {
	discounts, err := d.uc.GetAllExistingDiscounts(ctx, model.DiscountFilter{})
	if err != nil {
		return nil, err
	}

	var result []*svc.Discount
	for _, disc := range discounts {
		var productStrIDs []string
		for _, id := range disc.ApplicableProducts {
			productStrIDs = append(productStrIDs, id.Hex())
		}
		result = append(result, &svc.Discount{
			Id:                   disc.ID.Hex(),
			Name:                 disc.Name,
			Description:          disc.Description,
			DiscountPercentage:   float64(disc.DiscountPercentage),
			ApplicableProductIds: productStrIDs,
			IsActive:             disc.IsActive,
		})
	}

	return &svc.DiscountList{Discounts: result}, nil
}

func (d *Discount) GetProductsWithDiscounts(ctx context.Context, _ *emptypb.Empty) (*svc.ProductList, error) {
	products, err := d.uc.GetAllProductsWithDiscounts(ctx, model.DiscountFilter{})
	if err != nil {
		return nil, err
	}

	var protoProducts []*svc.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &svc.Product{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float64(p.Price),
			Stock:       int32(p.Stock),
			Category:    p.Category,
		})
	}
	return &svc.ProductList{Products: protoProducts}, nil
}

func (d *Discount) DeleteDiscount(ctx context.Context, req *svc.DiscountID) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	err = d.uc.DeleteDiscount(ctx, model.DiscountFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}