package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	StatusPending   ShipmentStatus = "pending"
	StatusPickedUp  ShipmentStatus = "picked_up"
	StatusInTransit ShipmentStatus = "in_transit"
	StatusDelivered ShipmentStatus = "delivered"
	StatusCancelled ShipmentStatus = "cancelled"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrShipmentNotFound        = errors.New("shipment not found")
	ErrInvalidShipment         = errors.New("invalid shipment data")
)

// Таблица разрешённых переходов между статусами
var allowedTransitions = map[ShipmentStatus][]ShipmentStatus{
	StatusPending:   {StatusPickedUp, StatusCancelled},
	StatusPickedUp:  {StatusInTransit, StatusCancelled},
	StatusInTransit: {StatusDelivered, StatusCancelled},
	StatusDelivered: {},
	StatusCancelled: {},
}

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	Status          ShipmentStatus
	DriverName      string
	DriverUnit      string
	Amount          float64
	DriverRevenue   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ShipmentEvent struct {
	ID         string
	ShipmentID string
	Status     ShipmentStatus
	Note       string
	CreatedAt  time.Time
}

func NewShipment(
	referenceNumber, origin, destination string,
	driverName, driverUnit string,
	amount, driverRevenue float64,
) (*Shipment, error) {
	if referenceNumber == "" {
		return nil, ErrInvalidShipment
	}
	if origin == "" || destination == "" {
		return nil, ErrInvalidShipment
	}
	if amount < 0 || driverRevenue < 0 {
		return nil, ErrInvalidShipment
	}

	now := time.Now().UTC()
	return &Shipment{
		ID:              uuid.NewString(),
		ReferenceNumber: referenceNumber,
		Origin:          origin,
		Destination:     destination,
		Status:          StatusPending,
		DriverName:      driverName,
		DriverUnit:      driverUnit,
		Amount:          amount,
		DriverRevenue:   driverRevenue,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (s *Shipment) CanTransitionTo(newStatus ShipmentStatus) bool {
	allowed, ok := allowedTransitions[s.Status]
	if !ok {
		return false
	}
	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}
	return false
}

func (s *Shipment) TransitionTo(newStatus ShipmentStatus) error {
	if !s.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}
	s.Status = newStatus
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func NewShipmentEvent(shipmentID string, status ShipmentStatus, note string) *ShipmentEvent {
	return &ShipmentEvent{
		ID:         uuid.NewString(),
		ShipmentID: shipmentID,
		Status:     status,
		Note:       note,
		CreatedAt:  time.Now().UTC(),
	}
}
