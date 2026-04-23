package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"shipping-app/internal/app/domain/events"
	"time"
)

type ReportDataFetcherPostgres struct {
	db               *sql.DB
	coordinatorEmail string
}

func NewReportDataFetcherPostgres(db *sql.DB, coordinatorEmail string) *ReportDataFetcherPostgres {
	return &ReportDataFetcherPostgres{
		db:               db,
		coordinatorEmail: coordinatorEmail,
	}
}

func (r *ReportDataFetcherPostgres) FetchReportData(packageId uint) (*events.ReportGeneratedEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			p.id,
			p.numpackage,
			p.status,
			ap.origin,
			ap.destination,
			s.id AS sender_id,
			s.name AS sender_name,
			s.email AS sender_email,
			id_info.id AS info_delivery_id,
			id_info.observations,
			id_info.reason_cancellation
		FROM packages p
		INNER JOIN addresspackages ap ON ap.id = p.idaddresspackage
		INNER JOIN senders s ON s.id = p.idsender
		LEFT JOIN informationdeliveries id_info ON id_info.idpackage = p.id
		WHERE p.id = $1
		LIMIT 1
	`
	var (
		evt                events.ReportGeneratedEvent
		infoDeliveryID     sql.NullInt64
		observation        sql.NullString
		reasonCancellation sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, packageId).Scan(
		&evt.PackageID,
		&evt.NumPackage,
		&evt.Status,
		&evt.Origin,
		&evt.Destination,
		&evt.SenderID,
		&evt.SenderName,
		&evt.SenderEmail,
		&infoDeliveryID,
		&observation,
		&reasonCancellation,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("paquete #%d no encontrado", packageId)
		}
		return nil, fmt.Errorf("fetch report data pkg#%d: %w ", packageId, err)
	}

	if infoDeliveryID.Valid {
		evt.InformationDeliveryID = uint(infoDeliveryID.Int64)
	}
	if observation.Valid {
		evt.Observation = &observation.String
	}
	if reasonCancellation.Valid {
		evt.ReasonCancellation = &reasonCancellation.String
	}

	evt.CoordinatorEmail = r.coordinatorEmail
	evt.OccurredAt = time.Now()

	return &evt, nil
}
