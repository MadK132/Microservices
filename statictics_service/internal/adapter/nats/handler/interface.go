// internal/adapter/nats/handler/interface.go
package handler

import (
	"context"
	"github.com/nats-io/nats.go"
)

type ProductUsecase interface {
	ProcessInventoryEvent(ctx context.Context, msg *nats.Msg) error
}

type OrderUsecase interface {
	ProcessOrderEvent(ctx context.Context, msg *nats.Msg) error
}
