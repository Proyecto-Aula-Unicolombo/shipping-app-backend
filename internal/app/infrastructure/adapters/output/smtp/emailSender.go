package smtp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"shipping-app/internal/app/domain/events"
	"shipping-app/internal/app/infrastructure/config"
)

type EmailSender struct {
	cfg      *config.SMTPConfig
	template map[string]*template.Template
}

//go:embed templates/*.html
var templateFiles embed.FS

func NewEmailSender(cfg *config.SMTPConfig) (*EmailSender, error) {
	fiels := map[string]string{
		"incident":    "templates/incident.html",
		"delivery":    "templates/delivery.html",
		"error_alert": "templates/errorAlert.html",
		"coordinator": "templates/coordinator.html",
	}

	parsed := make(map[string]*template.Template, len(fiels))
	for name, path := range fiels {
		t, err := template.ParseFS(templateFiles, path)
		if err != nil {
			return nil, fmt.Errorf("cargar template [%s]: %w", name, err)
		}
		parsed[name] = t
	}

	return &EmailSender{cfg: cfg, template: parsed}, nil
}

func (e *EmailSender) SendIncidentReport(evt *events.ReportGeneratedEvent) error {
	html, err := e.render("incident", evt)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Reporte de incidente - Paquete %s", evt.NumPackage)
	if err := e.sendMail(evt.SenderEmail, subject, html); err != nil {
		return err
	}

	return e.SendCoordinatorCopy(evt)
}

func (e *EmailSender) SendDeliveryReport(evt *events.ReportGeneratedEvent) error {
	html, err := e.render("delivery", evt)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("Reporte de Entegra - Paquete %s", evt.NumPackage)
	if err := e.sendMail(evt.SenderEmail, subject, html); err != nil {
		return err
	}

	return e.SendCoordinatorCopy(evt)
}

func (e *EmailSender) SendCoordinatorCopy(evt *events.ReportGeneratedEvent) error {
	html, err := e.render("coordinator", evt)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Report Paquete %s ---%s", evt.NumPackage, evt.Status)
	return e.sendMail(evt.CoordinatorEmail, subject, html)
}

func (e *EmailSender) SendErrorAlert(evt *events.ReportGeneratedEvent, reason string) error {
	type alertData struct {
		*events.ReportGeneratedEvent
		Reason string
	}

	html, err := e.render("error_alert", alertData{ReportGeneratedEvent: evt, Reason: reason})
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Error en reporte --- pkg#%d | %s", evt.PackageID, reason)
	return e.sendMail(evt.CoordinatorEmail, subject, html)
}

func (e *EmailSender) render(name string, data any) (string, error) {
	tmpl, ok := e.template[name]
	if !ok {
		return "", fmt.Errorf("template [%s] no registrado", name)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("ejecutar template [%s]: %w", name, err)
	}

	return buf.String(), nil
}

func (e *EmailSender) sendMail(to, subject, htmlBody string) error {
	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		e.cfg.From, to, subject, htmlBody,
	)

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	var auth smtp.Auth
	if e.cfg.User != "" {
		auth = smtp.PlainAuth("", e.cfg.User, e.cfg.Password, e.cfg.Host)
	}

	if err := smtp.SendMail(addr, auth, e.cfg.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp.SendMail a %s: %w", to, err)
	}

	log.Println("[SMTP] sent to %s | %s ", to, subject)

	return nil
}
