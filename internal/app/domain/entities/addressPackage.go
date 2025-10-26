package entities

type AddressPackage struct {
	ID                   uint
	Origin               string
	Destination          string
	DeliveryInstructions *string
}
