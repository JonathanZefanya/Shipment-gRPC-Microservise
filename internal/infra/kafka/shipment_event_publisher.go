package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"shipment/internal/domain"

	"github.com/segmentio/kafka-go"
)

type ShipmentEventPublisher struct {
	writer *kafka.Writer
}

type shipmentEventMessage struct {
	ID         string `json:"id"`
	ShipmentID string `json:"shipment_id"`
	Status     string `json:"status"`
	Note       string `json:"note"`
	CreatedAt  string `json:"created_at"`
}

func NewShipmentEventPublisher(brokers []string, topic string) *ShipmentEventPublisher {
	return &ShipmentEventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
	}
}

func (p *ShipmentEventPublisher) Publish(ctx context.Context, event *domain.ShipmentEvent) error {
	payload, err := json.Marshal(shipmentEventMessage{
		ID:         event.ID,
		ShipmentID: event.ShipmentID,
		Status:     string(event.Status),
		Note:       event.Note,
		CreatedAt:  event.CreatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
	})
	if err != nil {
		return fmt.Errorf("marshal shipment event: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ShipmentID),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("write shipment event: %w", err)
	}
	return nil
}

func (p *ShipmentEventPublisher) Close() error {
	return p.writer.Close()
}
