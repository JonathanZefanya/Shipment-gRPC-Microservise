package domain_test

import (
	"shipment/internal/domain"
	"testing"
)

func newTestShipment(t *testing.T) *domain.Shipment {
	t.Helper()
	s, err := domain.NewShipment(
		"REF-001", "Almaty", "Astana",
		"John Doe", "TRUCK-1",
		1000.0, 300.0,
	)
	if err != nil {
		t.Fatalf("unexpected error creating shipment: %v", err)
	}
	return s
}

func TestNewShipment_Success(t *testing.T) {
	s := newTestShipment(t)

	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if s.Status != domain.StatusPending {
		t.Errorf("expected status %s, got %s", domain.StatusPending, s.Status)
	}
	if s.ReferenceNumber != "REF-001" {
		t.Errorf("expected reference REF-001, got %s", s.ReferenceNumber)
	}
}

func TestNewShipment_EmptyReferenceNumber(t *testing.T) {
	_, err := domain.NewShipment(
		"", "Almaty", "Astana",
		"John Doe", "TRUCK-1",
		1000.0, 300.0,
	)
	if err == nil {
		t.Error("expected error for empty reference number")
	}
}

func TestNewShipment_EmptyOrigin(t *testing.T) {
	_, err := domain.NewShipment(
		"REF-001", "", "Astana",
		"John Doe", "TRUCK-1",
		1000.0, 300.0,
	)
	if err == nil {
		t.Error("expected error for empty origin")
	}
}

func TestNewShipment_NegativeAmount(t *testing.T) {
	_, err := domain.NewShipment(
		"REF-001", "Almaty", "Astana",
		"John Doe", "TRUCK-1",
		-100.0, 300.0,
	)
	if err == nil {
		t.Error("expected error for negative amount")
	}
}

func TestTransition_PendingToPickedUp(t *testing.T) {
	s := newTestShipment(t)

	if err := s.TransitionTo(domain.StatusPickedUp); err != nil {
		t.Errorf("expected valid transition, got error: %v", err)
	}
	if s.Status != domain.StatusPickedUp {
		t.Errorf("expected status %s, got %s", domain.StatusPickedUp, s.Status)
	}
}

func TestTransition_FullLifecycle(t *testing.T) {
	s := newTestShipment(t)

	transitions := []domain.ShipmentStatus{
		domain.StatusPickedUp,
		domain.StatusInTransit,
		domain.StatusDelivered,
	}

	for _, next := range transitions {
		if err := s.TransitionTo(next); err != nil {
			t.Errorf("expected valid transition to %s, got error: %v", next, err)
		}
	}

	if s.Status != domain.StatusDelivered {
		t.Errorf("expected final status %s, got %s", domain.StatusDelivered, s.Status)
	}
}

func TestTransition_PendingToCancelled(t *testing.T) {
	s := newTestShipment(t)

	if err := s.TransitionTo(domain.StatusCancelled); err != nil {
		t.Errorf("expected valid transition to cancelled, got error: %v", err)
	}
}

func TestTransition_PendingToDelivered_Invalid(t *testing.T) {
	s := newTestShipment(t)

	err := s.TransitionTo(domain.StatusDelivered)
	if err == nil {
		t.Error("expected error for invalid transition pending→delivered")
	}
}

func TestTransition_DeliveredToAny_Invalid(t *testing.T) {
	s := newTestShipment(t)

	_ = s.TransitionTo(domain.StatusPickedUp)
	_ = s.TransitionTo(domain.StatusInTransit)
	_ = s.TransitionTo(domain.StatusDelivered)

	invalidTargets := []domain.ShipmentStatus{
		domain.StatusPending,
		domain.StatusPickedUp,
		domain.StatusInTransit,
		domain.StatusCancelled,
	}

	for _, next := range invalidTargets {
		if err := s.TransitionTo(next); err == nil {
			t.Errorf("expected error for transition delivered→%s", next)
		}
	}
}

func TestTransition_CancelledToAny_Invalid(t *testing.T) {
	s := newTestShipment(t)
	_ = s.TransitionTo(domain.StatusCancelled)

	if err := s.TransitionTo(domain.StatusPending); err == nil {
		t.Error("expected error for transition cancelled→pending")
	}
}

func TestTransition_StatusUnchanged_AfterInvalid(t *testing.T) {
	s := newTestShipment(t)

	_ = s.TransitionTo(domain.StatusDelivered)

	if s.Status != domain.StatusPending {
		t.Errorf("status should remain %s after invalid transition, got %s",
			domain.StatusPending, s.Status)
	}
}
