// internal/usecase/stats.go
package usecase

import (
	"context"
	"log"
	"strings"
	_"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/proto-definitions/gen/inventory"
	"github.com/recktt77/proto-definitions/gen/orders"
	"github.com/recktt77/statistics_service/internal/repository"
	"google.golang.org/protobuf/proto"
)

type StatsUsecase struct {
	repo repository.StatsRepository
}

func NewStatsUsecase(repo repository.StatsRepository) *StatsUsecase {
	return &StatsUsecase{repo: repo}
}

func (s *StatsUsecase) ProcessInventoryEvent(ctx context.Context, msg *nats.Msg) error {
	var event inventory.InventoryEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		return err
	}

	log.Printf("Inventory event: %+v\n", event)
	s.repo.Save("inventory", strings.ToLower(event.Action), event.Timestamp)
	return nil
}

func (s *StatsUsecase) ProcessOrderEvent(ctx context.Context, msg *nats.Msg) error {
	var event order.OrderEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		return err
	}

	log.Printf("Order event: %+v\n", event)
	s.repo.Save("order", strings.ToLower(event.Action), event.Timestamp)
	return nil
}
