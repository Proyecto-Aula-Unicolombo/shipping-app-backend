package services

import (
	"fmt"
	"log"
	"shipping-app/internal/app/domain/events"
	"shipping-app/internal/app/domain/ports/output"
)

type ReportEmailService struct {
	notifier output.EmailNotifier
}

func NewReportEmailService(notifier output.EmailNotifier) *ReportEmailService {
	return &ReportEmailService{notifier: notifier}
}

func (s *ReportEmailService) HandleReportEvent(evt *events.ReportGeneratedEvent) error {
	log.Printf("[EMAIL_SERVICE] Procesando [%s] pkg#%d → %s",
		evt.EventType, evt.PackageID, evt.SenderEmail)

	switch evt.EventType {
	case events.EventIncidentReport:
		return s.notifier.SendIncidentReport(evt)
	case events.EventDeliveryReport:
		return s.notifier.SendDeliveryReport(evt)
	default:
		return fmt.Errorf("tipo de evento desconocido: %s", evt.EventType)
	}
}

func (s *ReportEmailService) HandleErrorAlert(evt *events.ErrorAlertEvent) error {
	log.Printf("[EMAIL_SERVICE] Alerta de error pkg#%d: %s", evt.PackageID, evt.Reason)

	reportEvt := &events.ReportGeneratedEvent{
		EventType:  events.EventErrorAlert,
		PackageID:  evt.PackageID,
		OccurredAt: evt.OccurredAt,
	}
	return s.notifier.SendErrorAlert(reportEvt, evt.Reason)
}
