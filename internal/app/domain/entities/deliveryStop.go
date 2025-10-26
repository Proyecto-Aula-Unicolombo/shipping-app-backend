package entities

import (
	"time"

	"github.com/twpayne/go-geom"
)

type DeliveryStop struct {
	ID           uint
	StopLocation *geom.Point
	TypeStop     string
	Timestamp    time.Time
	Description  *string
	Evidence     *string
	OrderID      uint
	Order        *Order
}
