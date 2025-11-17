package auth

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/externalServices/auth"
	"shipping-app/internal/utils"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	DriverID *uint  `json:"driver_id,omitempty"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyCredentials   = errors.New("email and password are required")
)

type LoginUseCase struct {
	userRepo   repository.UserRepository
	driverRepo repository.DriverRepository
	jwtService *auth.JWTService
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	driverRepo repository.DriverRepository,
	jwtService *auth.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		driverRepo: driverRepo,
		jwtService: jwtService,
	}
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	// Validar input
	if err := validateLoginInput(input); err != nil {
		return nil, err
	}

	// Buscar usuario por email
	user, err := uc.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar contraseña
	if !utils.VerifyPassword(input.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Si es conductor, obtener el driver ID
	var driverID *uint
	if user.Role == "driver" {
		driver, err := uc.driverRepo.GetDriverByUserID(user.ID)
		if err == nil && driver != nil {
			driverID = &driver.ID
		}
	}

	// Generar token JWT
	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.Role, driverID)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	// Preparar respuesta
	output := &LoginOutput{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Role:     user.Role,
			DriverID: driverID,
		},
	}

	return output, nil
}

func validateLoginInput(input LoginInput) error {
	if input.Email == "" || input.Password == "" {
		return ErrEmptyCredentials
	}
	return nil
}
