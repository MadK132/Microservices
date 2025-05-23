package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/proto-definitions/gen/inventory"
	"google.golang.org/protobuf/proto"
)

const PushTimeout = time.Second * 5

type InventoryProducer struct {
	nc *nats.Conn
}

func NewInventoryProducer(nc *nats.Conn) *InventoryProducer {
	return &InventoryProducer{nc: nc}
}

func (p *InventoryProducer) Push(ctx context.Context, productID string, action string) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	event := &inventory.InventoryEvent{
		ProductId: productID,
		Action:    action,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	subject := fmt.Sprintf("inventory.%s", action)

	if err := p.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("nats publish error: %w", err)
	}

	log.Println("Inventory event pushed:", event)
	return nil
}
