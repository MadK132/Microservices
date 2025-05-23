package backoffice

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/usecase"
	inventory "github.com/recktt77/proto-definitions/gen/inventory"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	inventory.UnimplementedProductServiceServer
	inventory.UnimplementedDiscountServiceServer

	cfg             *config.Server
	productUsecase  usecase.Product
	discountUsecase usecase.Discount
	grpcServer      *grpc.Server
}


func New(cfg *config.Server, productUC usecase.Product, discountUC usecase.Discount) Server {
	s := &server{
		
		cfg:             cfg,
		productUsecase:  productUC,
		discountUsecase: discountUC,
		grpcServer:      grpc.NewServer(),
	}

	inventory.RegisterProductServiceServer(s.grpcServer, s)
	inventory.RegisterDiscountServiceServer(s.grpcServer, s)

	reflection.Register(s.grpcServer)
	return s
}

func (s *server) Run(errCh chan<- error) {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(s.cfg.HTTPServer.Port))
	if err != nil {
		errCh <- err
		return
	}

	log.Println("gRPC server started on port", fmt.Sprint(s.cfg.HTTPServer.Port))
	if err := s.grpcServer.Serve(lis); err != nil {
		errCh <- err
	}
}

func (s *server) Stop() error {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
	return nil
}

//
// ProductService methods
//

func (s *server) CreateProduct(ctx context.Context, req *inventory.CreateProductRequest) (*inventory.Product, error) {
	product := model.Product{
		Name:  req.GetName(),
		Price: float64(req.GetPrice()),
	}

	created, err := s.productUsecase.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &inventory.Product{
		Id:    created.ID.Hex(),
		Name:  created.Name,
		Price: float64(created.Price),
	}, nil
}

func (s *server) GetProductByID(ctx context.Context, req *inventory.ProductID) (*inventory.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	product, err := s.productUsecase.GetByID(ctx, model.ProductFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	return &inventory.Product{
		Id:    product.ID.Hex(),
		Name:  product.Name,
		Price: float64(product.Price),
	}, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *inventory.UpdateProductRequest) (*inventory.Product, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	updated, err := s.productUsecase.Update(ctx, model.ProductUpdate{
		ID:    &id,
		Name:  req.Name,
		Price: ptrFloat64(req.GetPrice()),
	})
	if err != nil {
		return nil, err
	}
	return &inventory.Product{
		Id:    updated.ID.Hex(),
		Name:  updated.Name,
		Price: float64(updated.Price),
	}, nil
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func (s *server) DeleteProduct(ctx context.Context, req *inventory.ProductID) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	err = s.productUsecase.Delete(ctx, model.ProductFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *server) GetAllProducts(ctx context.Context, _ *emptypb.Empty) (*inventory.ProductList, error) {
	products, err := s.productUsecase.GetAll(ctx, model.ProductFilter{})
	if err != nil {
		return nil, err
	}
	var protoProducts []*inventory.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Price: float64(p.Price),
		})
	}
	return &inventory.ProductList{Products: protoProducts}, nil
}

func (s *server) GetProductsByIDs(ctx context.Context, req *inventory.ProductIDs) (*inventory.ProductList, error) {
	var ids []primitive.ObjectID
	for _, idStr := range req.GetIds() {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	products, err := s.productUsecase.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	var protoProducts []*inventory.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Price: float64(p.Price),
		})
	}
	return &inventory.ProductList{Products: protoProducts}, nil
}

// DiscountService methods
func (s *server) CreateDiscount(ctx context.Context, req *inventory.CreateDiscountRequest) (*inventory.Discount, error) {
	var productIDs []primitive.ObjectID
	for _, idStr := range req.GetApplicableProductIds() {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, err // или log.Println + continue, если хочешь игнорировать ошибки
		}
		productIDs = append(productIDs, id)
	}

	discount := model.Discount{
		Name:                 req.GetName(),
		Description:          req.GetDescription(),
		DiscountPercentage:   req.GetDiscountPercentage(),
		ApplicableProducts: productIDs,
		StartDate:            req.GetStartDate().AsTime(),
		EndDate:              req.GetEndDate().AsTime(),
	}
	
	_, err := s.discountUsecase.CreateDiscnout(ctx, discount)
	if err != nil {
		return nil, err
	}
	var stringIDs []string
	for _, objID := range discount.ApplicableProducts {
		stringIDs = append(stringIDs, objID.Hex())
	}

	return &inventory.Discount{
		Id:                   discount.ID.Hex(),
		Name:                 discount.Name,
		Description:          discount.Description,
		DiscountPercentage:   discount.DiscountPercentage,
		ApplicableProductIds: stringIDs,
		StartDate:            timestamppb.New(discount.StartDate),
		EndDate:              timestamppb.New(discount.EndDate),
		IsActive:             discount.IsActive,
	}, nil	
}

func (s *server) GetAllDiscounts(ctx context.Context, _ *emptypb.Empty) (*inventory.DiscountList, error) {
	discounts, err := s.discountUsecase.GetAllExistingDiscounts(ctx, model.DiscountFilter{})
	if err != nil {
		return nil, err
	}

	var protoDiscounts []*inventory.Discount
	for _, d := range discounts {
		var stringIDs []string
		for _, id := range d.ApplicableProducts {
			stringIDs = append(stringIDs, id.Hex())
		}

		protoDiscounts = append(protoDiscounts, &inventory.Discount{
			Id:                   d.ID.Hex(),
			Name:                 d.Name,
			Description:          d.Description,
			DiscountPercentage:   d.DiscountPercentage,
			ApplicableProductIds: stringIDs,
			StartDate:            timestamppb.New(d.StartDate),
			EndDate:              timestamppb.New(d.EndDate),
			IsActive:             d.IsActive,
		})
	}

	return &inventory.DiscountList{Discounts: protoDiscounts}, nil
}


func (s *server) GetProductsWithDiscounts(ctx context.Context, _ *emptypb.Empty) (*inventory.ProductList, error) {
	products, err := s.discountUsecase.GetAllProductsWithDiscounts(ctx, model.DiscountFilter{})
	if err != nil {
		return nil, err
	}
	var protoProducts []*inventory.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Price: float64(p.Price),
		})
	}
	return &inventory.ProductList{Products: protoProducts}, nil
}

func (s *server) DeleteDiscount(ctx context.Context, req *inventory.DiscountID) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	err = s.discountUsecase.DeleteDiscount(ctx, model.DiscountFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}


