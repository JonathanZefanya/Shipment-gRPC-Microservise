package domain

import "context"

type ShipmentRepository interface {
	Save(ctx context.Context, shipment *Shipment) error
	FindByID(ctx context.Context, id string) (*Shipment, error)
	Update(ctx context.Context, shipment *Shipment) error
}

type ShipmentEventRepository interface {
	Save(ctx context.Context, event *ShipmentEvent) error
	FindByShipmentID(ctx context.Context, shipmentID string) ([]*ShipmentEvent, error)
}

type ShipmentEventPublisher interface {
	Publish(ctx context.Context, event *ShipmentEvent) error
}
