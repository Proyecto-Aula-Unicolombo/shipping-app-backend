package usepackages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrToGetPackage = errors.New("error to get package")
	ErrToGetStatus  = errors.New("error to get package status")
	ErrCannotCancel = errors.New("package cannot be cancelled")
)

type CancelPackageUseCase struct {
	PackageRepo       repository.PackageRepository
	ComertialInfoRepo repository.ComercialInformationRepository
	ProviderTx        repository.TxProvider
}

func NewCancellPackageUseCase(
	packageRepo repository.PackageRepository,
	comertialInfoRepo repository.ComercialInformationRepository,
	providerTx repository.TxProvider,
) *CancelPackageUseCase {
	return &CancelPackageUseCase{
		PackageRepo:       packageRepo,
		ComertialInfoRepo: comertialInfoRepo,
		ProviderTx:        providerTx,
	}
}

func (u *CancelPackageUseCase) Execute(ctx context.Context, numPackage string) error {
	pkg, err := u.PackageRepo.GetByNumPackage(ctx, numPackage)
	if err != nil {
		log.Printf("ERROR: Failed to get package by numpackage %s: %v", numPackage, err)

		return ErrToGetPackage
	}

	pkgStatus, err := u.PackageRepo.GetStatusPackageToCancel(ctx, pkg.ID)
	if err != nil {
		return ErrToGetStatus
	}
	if pkgStatus.Status != "pendiente" {
		return ErrCannotCancel
	}

	tx, err := u.ProviderTx.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = u.ProviderTx.RollbackTx(ctx, tx)
		}
	}()

	if err := u.PackageRepo.DeletePackage(ctx, tx, pkg.ID); err != nil {
		return fmt.Errorf("delete package: %w", err)
	}

	if err := u.ComertialInfoRepo.Delete(ctx, tx, pkg.ComercialInformationID); err != nil {
		return fmt.Errorf("delete comercial information: %w", err)
	}

	if err := u.ProviderTx.CommitTx(ctx, tx); err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v", err)
		return fmt.Errorf("commit tx: %w", err)
	}
	commited = true

	return nil
}
