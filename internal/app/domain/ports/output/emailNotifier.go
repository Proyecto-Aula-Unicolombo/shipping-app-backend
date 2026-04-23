package output

import "shipping-app/internal/app/domain/events"

type EmailNotifier interface {
	SendIncidentReport(event *events.ReportGeneratedEvent) error
	SendDeliveryReport(event *events.ReportGeneratedEvent) error
	SendCoordinatorCopy(event *events.ReportGeneratedEvent) error
	SendErrorAlert(event *events.ReportGeneratedEvent, reason string) error
}
