package events

import "time"

const (
	EventIncidentReport = "report.incident"
	EventDeliveryReport = "report.delivery"
	EventErrorAlert     = "report.error"
)

// ReportGeneratedEvent es el evento que viaja por la cola de mensajes.
// Se construye con el JOIN de la DB y contiene todo lo necesario
// para que el Notification Worker envíe los emails sin hacer más consultas.
type ReportGeneratedEvent struct {
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`

	PackageID  uint   `json:"package_id"`
	NumPackage string `json:"num_package"`
	Status     string `json:"status"`

	Origin      string `json:"origin"`
	Destination string `json:"destination"`

	InformationDeliveryID uint    `json:"information_delivery_id"`
	Observation           *string `json:"observation,omitempty"`
	ReasonCancellation    *string `json:"reason_cancellation,omitempty"`
	PhotoDelivery         string  `json:"photo_delivery,omitempty"`
	SignatureReceived     *string `json:"signature_received,omitempty"`

	SenderID    uint   `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	SenderEmail string `json:"sender_email"`

	CoordinatorEmail string `json:"coordinator_email"`
}

// ErrorAlertEvent se publica cuando los datos del paquete no son válidos
// o cuando ocurre un error que impide generar el reporte.
type ErrorAlertEvent struct {
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
	PackageID  uint      `json:"package_id"`
	Reason     string    `json:"reason"`
}
