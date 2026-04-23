package input

import "shipping-app/internal/app/domain/events"

type ReportEventHandler interface {
	HandleReportEvent(event *events.ReportGeneratedEvent) error
	HandleErrorAlert(event *events.ErrorAlertEvent) error
}
