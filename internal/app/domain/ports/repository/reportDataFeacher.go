package repository

import "shipping-app/internal/app/domain/events"

type ReportDataFeatcher interface {
	FetchReportData(packageId uint) (*events.ReportGeneratedEvent, error)
}
