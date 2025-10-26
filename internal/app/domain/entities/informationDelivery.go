package entities

type InformationDelivery struct {
	ID                 uint
	Observation        *string
	SignatureReceived  *string
	PhotoDelivery      string
	ReasonCancellation *string
	PackageID          uint
	Package            *Package
}
