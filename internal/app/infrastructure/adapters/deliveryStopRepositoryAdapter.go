package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type DeliveryStopRepositoryPostgres struct {
	db *sql.DB
}

func NewDeliveryStopRepositoryPostgres(db *sql.DB) *DeliveryStopRepositoryPostgres {
	return &DeliveryStopRepositoryPostgres{db: db}
}

func (r *DeliveryStopRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, stop *entities.DeliveryStop) error {
	query := `
		INSERT INTO deliverystops (stoplocation, typestop, timestamp, description, evidence, idorder)
		VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, $4, $5, $6, $7)
		RETURNING id
	`

	lng := stop.StopLocation.X()
	lat := stop.StopLocation.Y()

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query,
			lng,
			lat,
			stop.TypeStop,
			stop.Timestamp,
			stop.Description,
			stop.Evidence,
			stop.OrderID,
		).Scan(&stop.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query,
			lng,
			lat,
			stop.TypeStop,
			stop.Timestamp,
			stop.Description,
			stop.Evidence,
			stop.OrderID,
		).Scan(&stop.ID)
	}

	if err != nil {
		log.Printf("ERROR creating delivery stop for order %d: %v", stop.OrderID, err)
		return fmt.Errorf("create delivery stop: %w", err)
	}

	log.Printf("✓ Delivery stop created: ID=%d, Type=%s, OrderID=%d", stop.ID, stop.TypeStop, stop.OrderID)
	return nil
}

func (r *DeliveryStopRepositoryPostgres) GetByID(ctx context.Context, id uint) (*entities.DeliveryStop, error) {
	query := `
		SELECT id, ST_AsBinary(stoplocation), typestop, timestamp, description, evidence, idorder
		FROM deliverystops
		WHERE id = $1
	`

	var stop entities.DeliveryStop
	var locationWKB []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stop.ID,
		&locationWKB,
		&stop.TypeStop,
		&stop.Timestamp,
		&stop.Description,
		&stop.Evidence,
		&stop.OrderID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrDeliveryStopNotFound
		}
		return nil, fmt.Errorf("get delivery stop by id: %w", err)
	}

	// Decodificar el punto WKB
	point, err := ewkb.Unmarshal(locationWKB)
	if err != nil {
		return nil, fmt.Errorf("unmarshal location: %w", err)
	}
	stop.StopLocation = point.(*geom.Point)

	return &stop, nil
}

func (r *DeliveryStopRepositoryPostgres) GetByOrderID(ctx context.Context, orderID uint) ([]*entities.DeliveryStop, error) {
	query := `
		SELECT id, ST_AsBinary(stoplocation), typestop, timestamp, description, evidence, idorder
		FROM deliverystops
		WHERE idorder = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get delivery stops by order id: %w", err)
	}
	defer rows.Close()

	var stops []*entities.DeliveryStop
	for rows.Next() {
		var stop entities.DeliveryStop
		var locationWKB []byte

		err := rows.Scan(
			&stop.ID,
			&locationWKB,
			&stop.TypeStop,
			&stop.Timestamp,
			&stop.Description,
			&stop.Evidence,
			&stop.OrderID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan delivery stop: %w", err)
		}

		// Decodificar el punto WKB
		point, err := ewkb.Unmarshal(locationWKB)
		if err != nil {
			return nil, fmt.Errorf("unmarshal location: %w", err)
		}
		stop.StopLocation = point.(*geom.Point)

		stops = append(stops, &stop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stops, nil
}

func (r *DeliveryStopRepositoryPostgres) ListByOrderIDWithLimit(ctx context.Context, orderID uint, limit int) ([]*entities.DeliveryStop, error) {
	query := `
		SELECT id, ST_AsBinary(stoplocation), typestop, timestamp, description, evidence, idorder
		FROM deliverystops
		WHERE idorder = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, orderID, limit)
	if err != nil {
		return nil, fmt.Errorf("list delivery stops by order id with limit: %w", err)
	}
	defer rows.Close()

	var stops []*entities.DeliveryStop
	for rows.Next() {
		var stop entities.DeliveryStop
		var locationWKB []byte

		err := rows.Scan(
			&stop.ID,
			&locationWKB,
			&stop.TypeStop,
			&stop.Timestamp,
			&stop.Description,
			&stop.Evidence,
			&stop.OrderID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan delivery stop: %w", err)
		}

		// Decodificar el punto WKB
		point, err := ewkb.Unmarshal(locationWKB)
		if err != nil {
			return nil, fmt.Errorf("unmarshal location: %w", err)
		}
		stop.StopLocation = point.(*geom.Point)

		stops = append(stops, &stop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stops, nil
}
