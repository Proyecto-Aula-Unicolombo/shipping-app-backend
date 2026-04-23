package output

import (
	"encoding/json"
	"fmt"
	"log"
	"shipping-app/internal/app/domain/events"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ReportsExchange    = "reports.exchange"
	QueueIncidents     = "reports.incidents"
	QueueDeliveries    = "reports.deliveries"
	QueueErrorAlerts   = "reports.erros"
	DeadletterExchange = "reports.dlx"
)

type EventPublisherRabbitMQ struct {
	channel *amqp.Channel
}

func NewEventPublisherRabbitMQ(ch *amqp.Channel) (*EventPublisherRabbitMQ, error) {
	publisher := &EventPublisherRabbitMQ{channel: ch}
	if err := publisher.setupTopology(); err != nil {
		return nil, fmt.Errorf("setup topologyÑ %w", &err)
	}

	return publisher, nil
}

func (p *EventPublisherRabbitMQ) setupTopology() error {
	// Exchange principal topic
	if err := p.channel.ExchangeDeclare(
		ReportsExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange %w", err)
	}

	// received message that fail
	if err := p.channel.ExchangeDeclare(
		DeadletterExchange,
		"fanout",
		true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("declare DLX: %w", err)
	}

	dqlArgs := amqp.Table{
		"x-dead-x-exchange": DeadletterExchange,
		"x-message-ttl":     int64(86400000), // 24h
	}

	queues := []struct {
		name       string
		routingKey string
	}{
		{QueueIncidents, events.EventIncidentReport},
		{QueueDeliveries, events.EventDeliveryReport},
		{QueueErrorAlerts, events.EventErrorAlert},
	}

	for _, queue := range queues {
		if _, err := p.channel.QueueDeclare(
			queue.name,
			true,
			false,
			false,
			false,
			dqlArgs,
		); err != nil {
			return fmt.Errorf("declare queue [%s]: %w", queue.name, err)
		}
		if err := p.channel.QueueBind(
			queue.name,
			queue.routingKey, // routing key = "report.incident" / "report.delivery" / "report.error"
			ReportsExchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("bind queue [%s]: %w", queue.name, err)
		}
	}

	log.Printf("[RABBITMQ] Exchange: %s | Colas: %s, %s, %s",
		ReportsExchange, QueueIncidents, QueueDeliveries, QueueErrorAlerts)
	return nil

}

func (p *EventPublisherRabbitMQ) publish(routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marashal payload: %w", err)
	}

	return p.channel.Publish(
		ReportsExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

func (p *EventPublisherRabbitMQ) PublishReportGenerated(evt *events.ReportGeneratedEvent) error {
	log.Printf("[RABBITMQ]  [%s] pkg#%d → sender: %s", evt.EventType, evt.PackageID, evt.SenderEmail)
	if err := p.publish(evt.EventType, evt); err != nil {
		return fmt.Errorf("publish report pkg#%d: %w", evt.PackageID, err)
	}
	return nil
}

func (p *EventPublisherRabbitMQ) PublishErrorAlert(evt *events.ErrorAlertEvent) error {
	log.Printf("[RABBITMQ] [error] pkg#%d: %s", evt.PackageID, evt.Reason)
	if err := p.publish(evt.EventType, evt); err != nil {
		return fmt.Errorf("publish error alert pkg#%d: %w", evt.PackageID, err)
	}

	return nil
}
