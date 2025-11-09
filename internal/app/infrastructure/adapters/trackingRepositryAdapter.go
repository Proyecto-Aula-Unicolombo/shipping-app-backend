package adapters

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"shipping-app/internal/app/domain/entities"

	"github.com/twpayne/go-geom/encoding/wkb"
)

type TrackingRepositoryAdapter struct {
	db *sql.DB
}

func NewTrackingRepositoryAdapter(db *sql.DB) *TrackingRepositoryAdapter {
	return &TrackingRepositoryAdapter{db: db}
}

func (r *TrackingRepositoryAdapter) RegisterTrack(ctx context.Context, track *entities.Track) error {
	insert := `INSERT INTO tracks (location, idorder) VALUES (ST_GeomFromWKB($1::bytea, 4326)::geography, $2) RETURNING id` // hacerle el cast a l tipo de dato

	locationWKB, err := wkb.Marshal(track.Location, binary.LittleEndian)
	if err != nil {
		return fmt.Errorf("failed to marshal location to WKB: %w", err)
	}

	err = r.db.QueryRowContext(ctx, insert, locationWKB, track.OrderID).Scan(&track.ID)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	return nil
}
