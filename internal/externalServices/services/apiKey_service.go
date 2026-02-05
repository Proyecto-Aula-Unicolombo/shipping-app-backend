package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
)

type APIKeyService struct {
	db *sql.DB
}

func NewAPIKeyService(db *sql.DB) *APIKeyService {
	return &APIKeyService{db: db}
}

// valida la API key y retorna el sender
func (s *APIKeyService) ValidateAPIKey(apiKey string) (*entities.Sender, error) {

	hashedKey := hashAPIKey(apiKey)

	query := `
		SELECT id, name, document, email, api_key, is_active
		FROM senders
		WHERE api_key = $1 AND is_active = true
	`

	var sender entities.Sender
	err := s.db.QueryRow(query, hashedKey).Scan(
		&sender.ID,
		&sender.Name,
		&sender.Document,
		&sender.Email,
		&sender.APIKey,
		&sender.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid API key")
		}
		return nil, fmt.Errorf("error validating API key: %w", err)
	}

	return &sender, nil
}

// crea un sender y le asigna una API Key
func (s *APIKeyService) CreateSenderWithAPIKey(name, document, address, phoneNumber, email string) (*entities.Sender, string, error) {
	// Generar API key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, "", err
	}

	hashedKey := hashAPIKey(apiKey)

	query := `
		INSERT INTO senders (name, document, address, phonenumber, email, api_key, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		RETURNING id
	`

	var sender entities.Sender
	sender.Name = name
	sender.Document = document
	sender.Email = email
	sender.APIKey = apiKey
	sender.IsActive = true

	err = s.db.QueryRow(query, name, document, address, phoneNumber, email, hashedKey).Scan(&sender.ID)
	if err != nil {
		return nil, "", fmt.Errorf("error creating sender: %w", err)
	}

	return &sender, apiKey, nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generando API key: %w", err)
	}

	plainKey := "sk_" + hex.EncodeToString(bytes)

	return plainKey, nil
}

func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}
