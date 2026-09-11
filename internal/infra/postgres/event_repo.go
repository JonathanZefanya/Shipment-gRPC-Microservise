package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"shipment/internal/domain"
)

type ShipmentEventRepo struct {
	db *sql.DB
}

func NewShipmentEventRepo(db *sql.DB) *ShipmentEventRepo {
	return &ShipmentEventRepo{db: db}
}

func (r *ShipmentEventRepo) Save(ctx context.Context, e *domain.ShipmentEvent) error {
	query := `
		INSERT INTO shipment_events (id, shipment_id, status, note, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.ShipmentID, string(e.Status), e.Note, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (r *ShipmentEventRepo) FindByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	query := `
		SELECT id, shipment_id, status, note, created_at
		FROM shipment_events
		WHERE shipment_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []*domain.ShipmentEvent
	for rows.Next() {
		var e domain.ShipmentEvent
		var status string
		if err := rows.Scan(&e.ID, &e.ShipmentID, &status, &e.Note, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		e.Status = domain.ShipmentStatus(status)
		events = append(events, &e)
	}
	return events, nil
}
