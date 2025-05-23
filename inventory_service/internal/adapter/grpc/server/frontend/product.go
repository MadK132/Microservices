package frontend

import (
	"context"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	svc "github.com/recktt77/proto-definitions/gen/inventory"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Product struct {
	svc.UnimplementedProductServiceServer

	uc ProductUsecase
}

func NewProduct(uc ProductUsecase) *Product{
	return &Product{
		uc:uc,
	}
}

func (p *Product) CreateProduct(ctx context.Context, req *svc.CreateProductRequest) (*svc.Product, error) {
	domainProduct := model.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       float64(req.GetPrice()),
		Stock:       int(req.GetStock()),
		Category:    req.GetCategory(),
	}

	created, err := p.uc.Create(ctx, domainProduct)
	if err != nil {
		return nil, err
	}

	return &svc.Product{
		Id:          created.ID.Hex(),
		Name:        created.Name,
		Description: created.Description,
		Price:       float64(created.Price),
		Stock:       int32(created.Stock),
		Category:    created.Category,
	}, nil
}

func (p *Product) UpdateProduct(ctx context.Context, req *svc.UpdateProductRequest) (*svc.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	var stock *int
	if req.Stock != nil {
		tmp := int(*req.Stock)
		stock = &tmp
	}


	update := model.ProductUpdate{
		ID:          &id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       stock,
		Category:    req.Category,
	}	


	updated, err := p.uc.Update(ctx, update)
	if err != nil {
		return nil, err
	}

	return &svc.Product{
		Id:          updated.ID.Hex(),
		Name:        updated.Name,
		Description: updated.Description,
		Price:       float64(updated.Price),
		Stock:       int32(updated.Stock),
		Category:    updated.Category,
	}, nil
}

func (p *Product) GetProductByID(ctx context.Context, req *svc.ProductID) (*svc.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	product, err := p.uc.GetByID(ctx, model.ProductFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	return &svc.Product{
		Id:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       float64(product.Price),
		Stock:       int32(product.Stock),
		Category:    product.Category,
	}, nil
}

func (p *Product) DeleteProduct(ctx context.Context, req *svc.ProductID) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	err = p.uc.Delete(ctx, model.ProductFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (p *Product) GetAllProducts(ctx context.Context, _ *emptypb.Empty) (*svc.ProductList, error) {
	products, err := p.uc.GetAll(ctx, model.ProductFilter{})
	if err != nil {
		return nil, err
	}

	var protoProducts []*svc.Product
	for _, prod := range products {
		protoProducts = append(protoProducts, &svc.Product{
			Id:          prod.ID.Hex(),
			Name:        prod.Name,
			Description: prod.Description,
			Price:       float64(prod.Price),
			Stock:       int32(prod.Stock),
			Category:    prod.Category,
		})
	}

	return &svc.ProductList{Products: protoProducts}, nil
}

func (p *Product) GetProductsByIDs(ctx context.Context, req *svc.ProductIDs) (*svc.ProductList, error) {
	var ids []primitive.ObjectID
	for _, idStr := range req.GetIds() {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	products, err := p.uc.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var protoProducts []*svc.Product
	for _, prod := range products {
		protoProducts = append(protoProducts, &svc.Product{
			Id:          prod.ID.Hex(),
			Name:        prod.Name,
			Description: prod.Description,
			Price:       float64(prod.Price),
			Stock:       int32(prod.Stock),
			Category:    prod.Category,
		})
	}

	return &svc.ProductList{Products: protoProducts}, nil
}
