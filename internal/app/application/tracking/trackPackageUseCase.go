package tracking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type TrackPackageInput struct {
	NumPackage *string
	PackageID  *uint
}

type LocationInfo struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

type TrackPackageResponse struct {
	PackageID       uint           `json:"package_id"`
	NumPackage      string         `json:"num_package"`
	Status          string         `json:"status"`
	Origin          string         `json:"origin"`
	Destination     string         `json:"destination"`
	CurrentLocation *LocationInfo  `json:"current_location"`
	ReceiverName    string         `json:"receiver_name"`
	ReceiverPhone   string         `json:"receiver_phone"`
	IsFragile       bool           `json:"is_fragile"`
	Weight          *float64       `json:"weight"`
	TrackingHistory []LocationInfo `json:"tracking_history"`
}

var (
	ErrInvalidTrackingInput = errors.New("must provide either numPackage or packageID")
	ErrUnauthorizedAccess   = errors.New("unauthorized: receiver does not have access to this package")
)

type TrackPackageUseCase struct {
	packageRepo repository.PackageRepository
	trackRepo   repository.TrackRepository
	orderRepo   repository.OrderRepository
	addressRepo repository.AddressPackageRepository
	reciverRepo repository.ReceiverRepository
}

func NewTrackPackageUseCase(
	packageRepo repository.PackageRepository,
	trackRepo repository.TrackRepository,
	orderRepo repository.OrderRepository,
	addressRepo repository.AddressPackageRepository,
	reciverRepo repository.ReceiverRepository,
) *TrackPackageUseCase {
	return &TrackPackageUseCase{
		packageRepo: packageRepo,
		trackRepo:   trackRepo,
		orderRepo:   orderRepo,
		addressRepo: addressRepo,
		reciverRepo: reciverRepo,
	}
}

func (uc *TrackPackageUseCase) Execute(ctx context.Context, input TrackPackageInput) (*TrackPackageResponse, error) {
	var pkg *entities.Package
	var err error

	// Obtener paquete por NumPackage o PackageID
	if input.NumPackage != nil && *input.NumPackage != "" {
		pkg, err = uc.packageRepo.GetByNumPackage(ctx, *input.NumPackage)
	} else if input.PackageID != nil {
		pkg, err = uc.packageRepo.GetByID(ctx, *input.PackageID)
	} else {
		return nil, ErrInvalidTrackingInput
	}

	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}

	// Construir respuesta base
	response := &TrackPackageResponse{
		PackageID:  pkg.ID,
		NumPackage: pkg.NumPackage,
		Status:     pkg.Status,
		IsFragile:  pkg.IsFragile,
		Weight:     pkg.Weight,
	}

	// Obtener información de dirección
	if pkg.AddressPackageID != 0 {
		address, err := uc.addressRepo.GetByID(ctx, pkg.AddressPackageID)
		if err == nil {
			response.Origin = address.Origin
			response.Destination = address.Destination
		}
	}

	// Obtener información del destinatario
	if pkg.ReceiverID != 0 {
		receiver, err := uc.reciverRepo.GetByID(ctx, pkg.ReceiverID)
		if err == nil {
			response.ReceiverName = receiver.Name + " " + receiver.LastName
		}
	}

	// Si el paquete tiene una orden asignada, obtener tracking
	if pkg.OrderID != nil {
		// Obtener historial de tracking
		tracks, err := uc.trackRepo.GetByOrderID(ctx, *pkg.OrderID)
		if err == nil && len(tracks) > 0 {
			// Agregar ubicación actual (la más reciente)
			latestTrack := tracks[0]
			response.CurrentLocation = &LocationInfo{
				Latitude:  latestTrack.Location.Y(),
				Longitude: latestTrack.Location.X(),
				Timestamp: latestTrack.Timestamp,
			}

			// Agregar historial de ubicaciones
			response.TrackingHistory = make([]LocationInfo, 0, len(tracks))
			for _, track := range tracks {
				response.TrackingHistory = append(response.TrackingHistory, LocationInfo{
					Latitude:  track.Location.Y(),
					Longitude: track.Location.X(),
					Timestamp: track.Timestamp,
				})
			}
		}
	}

	return response, nil
}
