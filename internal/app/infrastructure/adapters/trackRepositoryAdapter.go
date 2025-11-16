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

type TrackRepositoryPostgres struct {
	db *sql.DB
}

func NewTrackRepositoryPostgres(db *sql.DB) *TrackRepositoryPostgres {
	return &TrackRepositoryPostgres{db: db}
}

func (r *TrackRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, track *entities.Track) error {
	query := `
		INSERT INTO tracks (timestamp, location, idorder)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4)
		RETURNING id
	`

	lng := track.Location.X()
	lat := track.Location.Y()

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query,
			track.Timestamp,
			lng,
			lat,
			track.OrderID,
		).Scan(&track.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query,
			track.Timestamp,
			lng,
			lat,
			track.OrderID,
		).Scan(&track.ID)
	}

	if err != nil {
		log.Printf("ERROR creating track for order %d: %v", track.OrderID, err)
		return fmt.Errorf("create track: %w", err)
	}

	log.Printf("✓ Track created: ID=%d, OrderID=%d", track.ID, track.OrderID)
	return nil
}

func (r *TrackRepositoryPostgres) GetByOrderID(ctx context.Context, orderID uint) ([]*entities.Track, error) {
	query := `
		SELECT id, timestamp, ST_AsBinary(location), idorder
		FROM tracks
		WHERE idorder = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get tracks by order id: %w", err)
	}
	defer rows.Close()

	var tracks []*entities.Track
	for rows.Next() {
		var track entities.Track
		var locationWKB []byte

		err := rows.Scan(
			&track.ID,
			&track.Timestamp,
			&locationWKB,
			&track.OrderID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}

		// Decodificar el punto WKB
		point, err := ewkb.Unmarshal(locationWKB)
		if err != nil {
			return nil, fmt.Errorf("unmarshal location: %w", err)
		}
		track.Location = point.(*geom.Point)

		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tracks, nil
}

func (r *TrackRepositoryPostgres) GetLatestByOrderID(ctx context.Context, orderID uint) (*entities.Track, error) {
	query := `
		SELECT id, timestamp, ST_AsBinary(location), idorder
		FROM tracks
		WHERE idorder = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var track entities.Track
	var locationWKB []byte

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&track.ID,
		&track.Timestamp,
		&locationWKB,
		&track.OrderID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrTrackNotFound
		}
		return nil, fmt.Errorf("get latest track: %w", err)
	}

	// Decodificar el punto WKB
	point, err := ewkb.Unmarshal(locationWKB)
	if err != nil {
		return nil, fmt.Errorf("unmarshal location: %w", err)
	}
	track.Location = point.(*geom.Point)

	return &track, nil
}

func (r *TrackRepositoryPostgres) ListByOrderIDWithLimit(ctx context.Context, orderID uint, limit int) ([]*entities.Track, error) {
	query := `
		SELECT id, timestamp, ST_AsBinary(location), idorder
		FROM tracks
		WHERE idorder = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, orderID, limit)
	if err != nil {
		return nil, fmt.Errorf("list tracks by order id with limit: %w", err)
	}
	defer rows.Close()

	var tracks []*entities.Track
	for rows.Next() {
		var track entities.Track
		var locationWKB []byte

		err := rows.Scan(
			&track.ID,
			&track.Timestamp,
			&locationWKB,
			&track.OrderID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}

		// Decodificar el punto WKB
		point, err := ewkb.Unmarshal(locationWKB)
		if err != nil {
			return nil, fmt.Errorf("unmarshal location: %w", err)
		}
		track.Location = point.(*geom.Point)

		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tracks, nil
}
