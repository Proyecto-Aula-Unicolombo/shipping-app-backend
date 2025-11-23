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

func (r *DeliveryStopRepositoryPostgres) ListIncidents(ctx context.Context, status *string, driverID *uint, orderID *uint, limit int, offset int) ([]*entities.DeliveryStop, error) {
	query := `
		SELECT 
			ds.id, 
			ST_AsBinary(ds.stoplocation), 
			ds.typestop, 
			ds.timestamp, 
			ds.description, 
			ds.evidence, 
			ds.idorder,
			o.id,
			o.iddriver,
			o.status,
			d.id,
			u.id,
			u.name,
			u.lastname,
			u.email
		FROM deliverystops ds
		INNER JOIN orders o ON ds.idorder = o.id
		LEFT JOIN drivers d ON o.iddriver = d.id
		LEFT JOIN users u ON d.iduser = u.id
		WHERE ds.typestop = 'Incidente'
	`

	var args []interface{}
	argIndex := 1

	// Aplicar filtros opcionales
	if driverID != nil {
		query += fmt.Sprintf(" AND o.iddriver = $%d", argIndex)
		args = append(args, *driverID)
		argIndex++
	}

	if orderID != nil {
		query += fmt.Sprintf(" AND ds.idorder = $%d", argIndex)
		args = append(args, *orderID)
		argIndex++
	}

	query += " ORDER BY ds.timestamp DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list incidents: %w", err)
	}
	defer rows.Close()

	var stops []*entities.DeliveryStop
	for rows.Next() {
		var stop entities.DeliveryStop
		var locationWKB []byte

		// Order fields
		var orderIDVal uint
		var driverIDNullable *uint
		var orderStatus string

		// Driver fields (nullable)
		var driverIDVal sql.NullInt64
		var userID sql.NullInt64
		var userName sql.NullString
		var userLastName sql.NullString
		var userEmail sql.NullString

		err := rows.Scan(
			&stop.ID,
			&locationWKB,
			&stop.TypeStop,
			&stop.Timestamp,
			&stop.Description,
			&stop.Evidence,
			&stop.OrderID,
			// Order
			&orderIDVal,
			&driverIDNullable,
			&orderStatus,
			// Driver & User
			&driverIDVal,
			&userID,
			&userName,
			&userLastName,
			&userEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("scan incident: %w", err)
		}

		// Decodificar el punto WKB
		point, err := ewkb.Unmarshal(locationWKB)
		if err != nil {
			return nil, fmt.Errorf("unmarshal location: %w", err)
		}
		stop.StopLocation = point.(*geom.Point)

		// Construir la relación con Order y Driver
		stop.Order = &entities.Order{
			ID:       orderIDVal,
			Status:   orderStatus,
			DriverID: driverIDNullable,
		}

		if driverIDVal.Valid && userID.Valid {
			stop.Order.Driver = &entities.Driver{
				ID: uint(driverIDVal.Int64),
				User: &entities.User{
					ID:       uint(userID.Int64),
					Name:     userName.String,
					LastName: userLastName.String,
					Email:    userEmail.String,
				},
			}
		}

		stops = append(stops, &stop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stops, nil
}
