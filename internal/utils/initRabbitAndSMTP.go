package utils

import (
	"database/sql"
	"fmt"
	"log"
	"shipping-app/internal/app/domain/services"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/app/infrastructure/adapters/input/rabbitmq"
	adaptersOutput "shipping-app/internal/app/infrastructure/adapters/output"
	"shipping-app/internal/app/infrastructure/adapters/output/smtp"
	"shipping-app/internal/app/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Funtion to init Services RabbitMQ
func InitServicesRabbitMQ(cfg *config.Config, db *sql.DB) (*services.GenerateReportService, *amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	log.Printf("RABBITMQ URL: %v", cfg.RabbitMQ.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("error conectando RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error abriendo canal: %w", err)
	}
	ch.Qos(1, 0, false)

	// ── Adaptadores de salida ─────────────────────────────────
	generateReportSvc, reportEmailSvc, err := initServicesSMTP(cfg, ch, db)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Printf("error iniciando servicios SMTP: %v", err)
		return nil, nil, err
	}
	// ── Consumer (Notification Worker) ───────────────────────
	consumer := rabbitmq.NewReportConsumerRabbitMQ(ch, reportEmailSvc)
	if err := consumer.Start(); err != nil {
		ch.Close()
		conn.Close()
		log.Printf("error iniciando consumer: %v", err)
		return nil, nil, fmt.Errorf("error iniciando consumer: %w", err)
	}

	return generateReportSvc, conn, nil

}

// funtion to init components SMTP
func initServicesSMTP(cfg *config.Config, ch *amqp.Channel, db *sql.DB) (*services.GenerateReportService, *services.ReportEmailService, error) {
	// ── Adaptadores de salida ─────────────────────────────────
	eventPublisher, err := adaptersOutput.NewEventPublisherRabbitMQ(ch)
	if err != nil {
		return nil, nil, fmt.Errorf("error configurando RabbitMQ publisher: %v", err)
	}

	emailSender, err := smtp.NewEmailSender(&cfg.SMTP)
	if err != nil {
		return nil, nil, fmt.Errorf("error configurando SMTP: %v", err)
	}

	dataFetcher := adapters.NewReportDataFetcherPostgres(db, cfg.CoordinatorEmail)

	// ── Servicios de dominio ──────────────────────────────────
	generateReportSvc := services.NewGenerateReportsService(dataFetcher, eventPublisher)
	reportEmailSvc := services.NewReportEmailService(emailSender)

	return generateReportSvc, reportEmailSvc, nil
}
