package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	RabbitMQ         RabbitMQConfig
	SMTP             SMTPConfig
	CoordinatorEmail string
}

type RabbitMQConfig struct {
	URL string
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func Load() *Config {
	coordinatorEmail := getEnv("COORDINATOR_EMAIL", "")
	if coordinatorEmail == "" {
		log.Fatal("COORDINATOR_EMAIL es obligatorio — configúralo en las variables de entorno")
	}

	return &Config{
		CoordinatorEmail: coordinatorEmail,
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://admin:secret123@localhost:5672/"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnvInt("SMTP_PORT", 1025),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASS", ""),
			From:     getEnv("SMTP_FROM", "reportes@tuempresa.com"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Valor inválido para %s, usando default %d", key, fallback)
			return fallback
		}
		return i
	}
	return fallback
}
