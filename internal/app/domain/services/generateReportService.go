package services

import (
	"fmt"
	"log"
	"shipping-app/internal/app/domain/events"
	"shipping-app/internal/app/domain/ports/output"
	"shipping-app/internal/app/domain/ports/repository"
	"strings"
	"time"
)

type GenerateReportService struct {
	dataFetcher    repository.ReportDataFeatcher
	eventPublisher output.EventPublisher
}

func NewGenerateReportsService(
	dataFetcher repository.ReportDataFeatcher,
	eventPublisher output.EventPublisher,
) *GenerateReportService {
	return &GenerateReportService{
		dataFetcher:    dataFetcher,
		eventPublisher: eventPublisher,
	}
}

func (s *GenerateReportService) Execute(packageId uint) error {
	log.Printf("[GENERATE_REPORT] consulting data of package  #%d", packageId)

	evt, err := s.dataFetcher.FetchReportData(packageId)
	if err != nil {
		log.Printf("[GENERATE_REPORT] Error to get data pkg#%d: %v", packageId, err)
		s.publishErrorAlert(packageId, fmt.Sprintf("error obteniendo datos: %v", err))
		return fmt.Errorf("fetch report data: %w", err)
	}

	if err := s.validate(evt); err != nil {
		log.Printf("[GENERATE REPORT] Invalid data pkg#%d: %v", packageId, err)
		s.publishErrorAlert(packageId, err.Error())
		return err
	}

	evt.EventType = s.resolveEventype(evt.Status)
	evt.OccurredAt = time.Now()

	log.Printf("[GENERATE REPORT] Event is publishing [%s] pkg#%d → %s",
		evt.EventType, packageId, evt.SenderEmail)

	if err := s.eventPublisher.PublishReportGenerated(evt); err != nil {
		s.publishErrorAlert(packageId, fmt.Sprintf("error publicando evento: %v", err))
		return fmt.Errorf("publish report: %w", err)
	}

	log.Printf("GENERATE REPORT event is publish correctly pkg#%d", packageId)

	return nil
}

func (s *GenerateReportService) validate(evt *events.ReportGeneratedEvent) error {
	if evt.SenderEmail == "" {
		return fmt.Errorf("sender email vacío para paquete #%d", evt.PackageID)
	}

	if evt.CoordinatorEmail == "" {
		return fmt.Errorf("coordinator email no configurado")
	}

	if evt.InformationDeliveryID == 0 {
		return fmt.Errorf("informationdelivery no existe para el paquete #%d", evt.PackageID)
	}

	if evt.NumPackage == "" {
		return fmt.Errorf("numpackage vacío para paquete #%d", evt.PackageID)
	}

	return nil
}

func (s *GenerateReportService) resolveEventype(ststus string) string {
	switch strings.ToLower(ststus) {
	case "incidente":
		return events.EventIncidentReport
	case "entregado":
		return events.EventDeliveryReport
	default:
		return events.EventDeliveryReport
	}
}

func (s *GenerateReportService) publishErrorAlert(packageId uint, reason string) {
	alert := &events.ErrorAlertEvent{
		EventType:  events.EventErrorAlert,
		OccurredAt: time.Now(),
		PackageID:  packageId,
		Reason:     reason,
	}

	if err := s.eventPublisher.PublishErrorAlert(alert); err != nil {
		log.Printf("[GENERATE REPORT] Can not pubish alert pkg#%d: %v", packageId, err)
	}
}
