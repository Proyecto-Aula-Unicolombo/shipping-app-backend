package entities

import "time"

type Package struct {
	ID                     uint
	NumPackage             string
	Status                 string
	DescriptionContent     *string
	Weight                 *float64
	Dimension              *string
	DeclaredValue          *float64
	TypePackage            *string
	IsFragile              bool
	CreatedAt              time.Time
	UpdatedAt              *time.Time
	AddressPackageID       uint
	ComercialInformationID uint
	SenderID               uint
	ReceiverID             uint
	OrderID                *uint
	AddressPackage         *AddressPackage
	ComercialInformation   *ComercialInformation
	Sender                 *Sender
	Receiver               *Receiver
	Order                  *Order
}
