package entities

import "time"

type StatusDelivery struct {
	ID                    uint
	Status                string
	Priority              string
	DateEstimatedDelivery *time.Time
	DateRealDelivery      *time.Time
}
