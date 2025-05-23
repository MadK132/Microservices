package server

import (
	"context"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/grpc/server/frontend"
	frontendsvc "github.com/recktt77/proto-definitions/gen/inventory"
	"github.com/recktt77/Microservices-First-/inventory_service/config"

)

type API struct {
	s             *grpc.Server
	cfg           config.GRPCServer
	addr          string
	productUsecase ProductUsecase
	discountUsecase DiscountUsecase
	reviewUsecase ReviewUsecase
}

func New(
	cfg config.GRPCServer,
	productUsecase ProductUsecase,
	discountUsecase DiscountUsecase,
	reviewUsecase ReviewUsecase,
) *API {
	return &API{
		cfg:           cfg,
		addr:          fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		productUsecase: productUsecase,
		discountUsecase: discountUsecase,
		reviewUsecase: reviewUsecase,
	}
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
	go func() {
		log.Println(ctx, "gRPC server starting listen", fmt.Sprintf("addr: %s", a.addr))

		if err := a.run(ctx); err != nil {
			errCh <- fmt.Errorf("can't start grpc server: %w", err)

			return
		}
	}()
}

func (a *API) Stop(ctx context.Context) error {
	if a.s == nil {
		return nil
	}

	stopped := make(chan struct{})
	go func() {
		a.s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		a.s.Stop()
	case <-stopped:
	}

	return nil
}

func (a *API) run(ctx context.Context) error {
	a.s = grpc.NewServer(a.setOptions(ctx)...)

	frontendsvc.RegisterProductServiceServer(a.s, frontend.NewProduct(a.productUsecase))
	frontendsvc.RegisterDiscountServiceServer(a.s, frontend.NewDiscount(a.discountUsecase))
	frontendsvc.RegisterReviewServiceServer(a.s, frontend.NewReview(a.reviewUsecase))

	reflection.Register(a.s)

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	err = a.s.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}