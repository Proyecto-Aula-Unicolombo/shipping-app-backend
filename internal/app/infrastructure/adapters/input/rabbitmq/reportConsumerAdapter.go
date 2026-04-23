package rabbitmq

import (
	"encoding/json"
	"log"
	"shipping-app/internal/app/domain/events"
	inputPorts "shipping-app/internal/app/domain/ports/input"
	"shipping-app/internal/app/infrastructure/adapters/output"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReportConsumeRabbitMQ struct {
	channel *amqp.Channel
	handler inputPorts.ReportEventHandler
}

func NewReportConsumerRabbitMQ(ch *amqp.Channel, handler inputPorts.ReportEventHandler) *ReportConsumeRabbitMQ {
	return &ReportConsumeRabbitMQ{
		channel: ch,
		handler: handler,
	}
}

func (c *ReportConsumeRabbitMQ) Start() error {
	queues := []struct {
		name    string
		process func(amqp.Delivery)
	}{
		{output.QueueIncidents, c.processReportEvent},
		{output.QueueDeliveries, c.processReportEvent},
		{output.QueueErrorAlerts, c.processErrorAlert},
	}

	for _, queue := range queues {
		msgs, err := c.channel.Consume(
			queue.name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
		log.Printf("[CONSUMER] escuchando cola: %s", queue.name)

		processFn := queue.process
		go func() {
			for msg := range msgs {
				processFn(msg)
			}
		}()
	}
	return nil
}

func (c *ReportConsumeRabbitMQ) processReportEvent(msg amqp.Delivery) {
	var evt events.ReportGeneratedEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("[CONSUMER] error %v", err)
		msg.Nack(false, false) // send to DLX
		return
	}

	log.Printf("[CONSUMER] [%s] pkg#%d", evt.EventType, evt.PackageID)
	if err := c.handler.HandleReportEvent(&evt); err != nil {
		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
}

func (c *ReportConsumeRabbitMQ) processErrorAlert(msg amqp.Delivery) {
	var evt events.ErrorAlertEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("[CONSUMER]  Error deserializando ErrorAlertEvent: %v", err)
		msg.Nack(false, false)
		return
	}
	log.Printf("[CONSUMER] Alerta de error pkg#%d: %s", evt.PackageID, evt.Reason)

	if err := c.handler.HandleErrorAlert(&evt); err != nil {
		log.Printf("[CONSUMER] Error procesando alerta pkg#%d: %v", evt.PackageID, err)
		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
}
