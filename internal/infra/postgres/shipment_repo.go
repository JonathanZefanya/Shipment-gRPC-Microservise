package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shipment/internal/domain"

	_ "github.com/lib/pq"
)

type ShipmentRepo struct {
	db *sql.DB
}

func NewShipmentRepo(db *sql.DB) *ShipmentRepo {
	return &ShipmentRepo{db: db}
}

func (r *ShipmentRepo) Save(ctx context.Context, s *domain.Shipment) error {
	query := `
		INSERT INTO shipments 
			(id, reference_number, origin, destination, status,
			 driver_name, driver_unit, amount, driver_revenue, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.ReferenceNumber, s.Origin, s.Destination, string(s.Status),
		s.DriverName, s.DriverUnit, s.Amount, s.DriverRevenue,
		s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert shipment: %w", err)
	}
	return nil
}

func (r *ShipmentRepo) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	query := `
		SELECT id, reference_number, origin, destination, status,
		       driver_name, driver_unit, amount, driver_revenue, created_at, updated_at
		FROM shipments WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.Shipment
	var status string

	err := row.Scan(
		&s.ID, &s.ReferenceNumber, &s.Origin, &s.Destination, &status,
		&s.DriverName, &s.DriverUnit, &s.Amount, &s.DriverRevenue,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrShipmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan shipment: %w", err)
	}

	s.Status = domain.ShipmentStatus(status)
	return &s, nil
}

func (r *ShipmentRepo) Update(ctx context.Context, s *domain.Shipment) error {
	query := `
		UPDATE shipments 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, string(s.Status), s.UpdatedAt, s.ID)
	if err != nil {
		return fmt.Errorf("update shipment: %w", err)
	}
	return nil
}
