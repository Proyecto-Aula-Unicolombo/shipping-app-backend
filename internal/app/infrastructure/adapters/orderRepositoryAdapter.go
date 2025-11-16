package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
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
			iddriver,
			idvehicle
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	args := []interface{}{
		order.CreateAt,
		order.AssignedAt,
		order.Observation,
		order.Status,
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
		SET assigned_at = $1, observation = $2, status = $3, iddriver = $4, idvehicle = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		order.AssignedAt,
		order.Observation,
		order.Status,
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
		SELECT id, create_at, assigned_at, observation, status, iddriver, idvehicle
		FROM orders
		WHERE id = $1
	`

	var order entities.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.CreateAt,
		&order.AssignedAt,
		&order.Observation,
		&order.Status,
		&order.DriverID,
		&order.VehicleID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	return &order, nil
}

func (r *OrderRepositoryPostgres) List(ctx context.Context, limit, offset int) ([]*entities.Order, error) {
	query := `
		SELECT id, create_at, assigned_at, observation, status, iddriver, idvehicle
		FROM orders
		ORDER BY create_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
		SELECT id, create_at, assigned_at, observation, status, iddriver, idvehicle
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
