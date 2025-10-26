package entities

import (
	"time"

	"github.com/twpayne/go-geom"
)

type Track struct {
	ID        uint
	Timestamp time.Time
	Location  *geom.Point
	OrderID   uint
	Order     *Order
}
