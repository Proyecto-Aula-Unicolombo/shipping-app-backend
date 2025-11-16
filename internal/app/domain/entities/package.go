package entities

import "time"

type Package struct {
	ID                     uint
	NumPackage             string
	StartStatus            string
	DescriptionContent     *string
	Weight                 *float64
	Dimension              *string
	DeclaredValue          *float64
	TypePackage            *string
	IsFragile              bool
	CreatedAt              time.Time
	UpdatedAt              *time.Time
	AddressPackageID       uint
	StatusDeliveryID       uint
	ComercialInformationID uint
	SenderID               uint
	ReceiverID             uint
	OrderID                *uint
	AddressPackage         *AddressPackage
	StatusDelivery         *StatusDelivery
	ComercialInformation   *ComercialInformation
	Sender                 *Sender
	Receiver               *Receiver
	Order                  *Order
}
