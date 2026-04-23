package output

import "shipping-app/internal/app/domain/events"

type EventPublisher interface {
	PublishReportGenerated(event *events.ReportGeneratedEvent) error
	PublishErrorAlert(event *events.ErrorAlertEvent) error
}
