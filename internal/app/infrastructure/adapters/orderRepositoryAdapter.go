package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/lib/pq"
)

type OrderRepositoryPostgres struct {
	db *sql.DB
}

func NewOrderRepositoryPostgres(db *sql.DB) *OrderRepositoryPostgres {
	return &OrderRepositoryPostgres{db: db}
}

func (r *OrderRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, order *entities.Order) error {
	query := `
		INSERT INTO orders (
			create_at,
			assigned_at,
			observation,
			status,
			typeservice,
			iddriver,
			idvehicle
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	args := []interface{}{
		order.CreateAt,
		order.AssignedAt,
		order.Observation,
		order.Status,
		order.TypeService,
		order.DriverID,
		order.VehicleID,
	}

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&order.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&order.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("order create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("order create timeout: %w", err)
		}
		log.Printf("ERROR creating order: %v", err)
		return fmt.Errorf("order create: %w", err)
	}

	log.Printf("✓ Order created successfully: ID=%d, Status=%s", order.ID, order.Status)
	return nil
}

func (r *OrderRepositoryPostgres) Update(ctx context.Context, order *entities.Order) error {
	query := `
		UPDATE orders 
		SET assigned_at = $1, observation = $2, status = $3, typeservice = $4, iddriver = $5, idvehicle = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		order.AssignedAt,
		order.Observation,
		order.Status,
		order.TypeService,
		order.DriverID,
		order.VehicleID,
		order.ID,
	)

	if err != nil {
		log.Printf("ERROR updating order %d: %v", order.ID, err)
		return fmt.Errorf("update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrOrderNotFound
	}

	log.Printf("✓ Order updated successfully: ID=%d", order.ID)
	return nil
}

func (r *OrderRepositoryPostgres) UpdateStatus(ctx context.Context, id uint, status string, observation *string) error {
	query := `
		UPDATE orders 
		SET status = $1, observation = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, observation, id)
	if err != nil {
		log.Printf("ERROR updating order status %d: %v", id, err)
		return fmt.Errorf("update order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrOrderNotFound
	}

	log.Printf("✓ Order status updated: ID=%d, Status=%s", id, status)
	return nil
}

func (r *OrderRepositoryPostgres) AssignDriverAndVehicle(ctx context.Context, id uint, driverID, vehicleID uint) error {
	query := `
		UPDATE orders 
		SET iddriver = $1, idvehicle = $2, assigned_at = CURRENT_TIMESTAMP, status = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, driverID, vehicleID, "En camino", id)
	if err != nil {
		log.Printf("ERROR assigning driver/vehicle to order %d: %v", id, err)
		return fmt.Errorf("assign driver and vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrOrderNotFound
	}

	log.Printf("✓ Order assigned: ID=%d, DriverID=%d, VehicleID=%d", id, driverID, vehicleID)
	return nil
}

func (r *OrderRepositoryPostgres) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("ERROR deleting order %d: %v", id, err)
		return fmt.Errorf("delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrOrderNotFound
	}

	log.Printf("✓ Order deleted: ID=%d", id)
	return nil
}

func (r *OrderRepositoryPostgres) DeleteWithTx(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `DELETE FROM orders WHERE id = $1`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, id)
	} else {
		result, err = r.db.ExecContext(ctx, query, id)
	}

	if err != nil {
		log.Printf("ERROR deleting order %d: %v", id, err)
		return fmt.Errorf("delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrOrderNotFound
	}

	log.Printf("✓ Order deleted: ID=%d", id)
	return nil
}

func (r *OrderRepositoryPostgres) GetByID(ctx context.Context, id uint) (*entities.Order, error) {
	query := `
		SELECT 
			o.id, 
			o.create_at, 
			o.assigned_at, 
			o.observation, 
			o.status, 
			o.typeservice, 
			o.iddriver, 
			o.idvehicle,
			COALESCE(array_agg(p.id) FILTER (WHERE p.id IS NOT NULL), '{}') as package_ids
		FROM orders o
		LEFT JOIN packages p ON o.id = p.idorder
		WHERE o.id = $1
		GROUP BY o.id, o.create_at, o.assigned_at, o.observation, o.status, o.typeservice, o.iddriver, o.idvehicle
	`

	var order entities.Order
	var packageIDs []int64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.CreateAt,
		&order.AssignedAt,
		&order.Observation,
		&order.Status,
		&order.TypeService,
		&order.DriverID,
		&order.VehicleID,
		(*pq.Int64Array)(&packageIDs),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	// Convert []int64 to []uint
	order.PackageIDs = make([]uint, len(packageIDs))
	for i, id := range packageIDs {
		order.PackageIDs[i] = uint(id)
	}

	return &order, nil
}

func (r *OrderRepositoryPostgres) List(ctx context.Context, orderID uint, limit, offset int, typeService, status string) ([]*entities.Order, error) {
	query := `
		SELECT id, create_at, assigned_at, observation, status, typeservice, iddriver, idvehicle
		FROM orders
		WHERE 1=1
	`

	args := []interface{}{}
	argPosition := 1

	if typeService != "" {
		query += fmt.Sprintf(" AND typeservice = $%d", argPosition)
		args = append(args, typeService)
		argPosition++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPosition)
		args = append(args, status)
		argPosition++
	}

	if orderID != 0 {
		query += fmt.Sprintf(" AND id = $%d", argPosition)
		args = append(args, orderID)
		argPosition++
	}

	query += fmt.Sprintf(" ORDER BY create_at DESC LIMIT $%d OFFSET $%d", argPosition, argPosition+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []*entities.Order
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(
			&order.ID,
			&order.CreateAt,
			&order.AssignedAt,
			&order.Observation,
			&order.Status,
			&order.TypeService,
			&order.DriverID,
			&order.VehicleID,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) ListByDriver(ctx context.Context, driverID uint, limit, offset int) ([]*entities.Order, error) {
	query := `
		SELECT id, create_at, assigned_at, observation, status, typeservice, iddriver, idvehicle
		FROM orders
		WHERE iddriver = $1
		ORDER BY create_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, driverID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list orders by driver: %w", err)
	}
	defer rows.Close()

	var orders []*entities.Order
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(
			&order.ID,
			&order.CreateAt,
			&order.AssignedAt,
			&order.Observation,
			&order.Status,
			&order.TypeService,
			&order.DriverID,
			&order.VehicleID,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) Count(ctx context.Context, typeService, status string) (int64, error) {
	query := `
		SELECT COUNT(*) FROM orders WHERE typeservice ILIKE $1 AND status ILIKE $2
	`
	var count int64
	err := r.db.QueryRowContext(ctx, query, "%"+typeService+"%", "%"+status+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count orders: %w", err)
	}

	return count, nil
}

func (r *OrderRepositoryPostgres) CountByDriver(ctx context.Context, driverID uint) (int64, error) {
	query := `
		SELECT COUNT(*) FROM orders WHERE iddriver = $1
		`
	var count int64
	err := r.db.QueryRowContext(ctx, query, driverID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count orders by driver: %w", err)
	}
	return count, nil
}
