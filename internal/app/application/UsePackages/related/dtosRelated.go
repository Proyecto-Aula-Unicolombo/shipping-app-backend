package usepackages

import "time"

type AdressPackageInput struct {
	Origin               string  `json:"origin"`
	Destination          string  `json:"destination"`
	DeliveryInstructions *string `json:"delivery_instructions"`
}

type StatusDeliveryInput struct {
	Status                string     `json:"status"`
	Priority              string     `json:"priority"`
	DateEstimatedDelivery *time.Time `json:"date_estimated_delivery"`
	DateRealDelivery      *time.Time `json:"date_real_delivery"`
}

type ComercialInformationInput struct {
	CostSending float64 `json:"cost_sending"`
	IsPaid      bool    `json:"is_paid"`
}

type SenderInput struct {
	Name        string `json:"name"`
	Document    string `json:"document"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type ReceiverInput struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type AdressPackageResponse struct {
	Origin               string  `json:"origin"`
	Destination          string  `json:"destination"`
	DeliveryInstructions *string `json:"delivery_instructions"`
}

type StatusDeliveryResponse struct {
	Status                string     `json:"status"`
	Priority              string     `json:"priority"`
	DateEstimatedDelivery *time.Time `json:"date_estimated_delivery"`
	DateRealDelivery      *time.Time `json:"date_real_delivery"`
}

type ComercialInformationResponse struct {
	CostSending float64 `json:"cost_sending"`
	IsPaid      bool    `json:"is_paid"`
}

type SenderResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReceiverResponse struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}
