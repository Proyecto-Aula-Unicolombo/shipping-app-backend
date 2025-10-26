package entities

import "time"

type Package struct {
	ID                     uint
	NumPackage             int64
	StartStatus            string
	DescriptionContent     *string
	Weight                 *float64
	Dimension              *float64
	DeclaredValue          *float64
	TypePackage            *string
	IsFragile              bool
	CreatedAt              time.Time
	UpdatedAt              *time.Time
	AddressPackageID       uint
	StatusDeliveryID       uint
	InformationDeliveryID  *uint
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
