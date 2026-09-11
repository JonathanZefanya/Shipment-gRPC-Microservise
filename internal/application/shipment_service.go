package application

import (
	"context"
	"fmt"
	"shipment/internal/domain"
)

type ShipmentService struct {
	shipments domain.ShipmentRepository
	events    domain.ShipmentEventRepository
	publisher domain.ShipmentEventPublisher
}

func NewShipmentService(shipments domain.ShipmentRepository, events domain.ShipmentEventRepository, publisher domain.ShipmentEventPublisher) *ShipmentService {
	return &ShipmentService{
		shipments: shipments,
		events:    events,
		publisher: publisher,
	}
}

func (s *ShipmentService) CreateShipment(ctx context.Context, referenceNumber, origin, destination string, driverName, driverUnit string, amount, driverRevenue float64) (*domain.Shipment, error) {
	shipment, err := domain.NewShipment(
		referenceNumber, origin, destination,
		driverName, driverUnit,
		amount, driverRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("create shipment: %w", err)
	}

	if err := s.shipments.Save(ctx, shipment); err != nil {
		return nil, fmt.Errorf("save shipment: %w", err)
	}

	event := domain.NewShipmentEvent(shipment.ID, domain.StatusPending, "shipment created")
	if err := s.events.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save initial event: %w", err)
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return nil, fmt.Errorf("publish initial event: %w", err)
	}

	return shipment, nil
}

// GetShipment — shipment by ID
func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	shipment, err := s.shipments.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get shipment: %w", err)
	}
	return shipment, nil
}

func (s *ShipmentService) AddShipmentEvent(ctx context.Context, shipmentID string, newStatus domain.ShipmentStatus, note string) (*domain.ShipmentEvent, error) {
	shipment, err := s.shipments.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("find shipment: %w", err)
	}

	if err := shipment.TransitionTo(newStatus); err != nil {
		return nil, fmt.Errorf("transition status: %w", err)
	}

	if err := s.shipments.Update(ctx, shipment); err != nil {
		return nil, fmt.Errorf("update shipment: %w", err)
	}

	event := domain.NewShipmentEvent(shipmentID, newStatus, note)
	if err := s.events.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}

	return event, nil
}

func (s *ShipmentService) GetShipmentHistory(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	if _, err := s.shipments.FindByID(ctx, shipmentID); err != nil {
		return nil, fmt.Errorf("find shipment: %w", err)
	}

	events, err := s.events.FindByShipmentID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	return events, nil
}
