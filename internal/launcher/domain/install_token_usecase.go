package domain

import (
	"context"
	"pb_launcher/internal/launcher/domain/repositories"
)

type CleanServiceInstallTokenUsecase struct {
	repository repositories.ServiceRepository
}

func NewCleanServiceInstallTokenUsecase(repository repositories.ServiceRepository) *CleanServiceInstallTokenUsecase {
	return &CleanServiceInstallTokenUsecase{repository: repository}
}

func (u *CleanServiceInstallTokenUsecase) SetInstallToken(ctx context.Context, serviceID, token string) error {
	if token == "" {
		return nil
	}
	return u.repository.SetServiceInstallToken(ctx, serviceID, token)
}

func (u *CleanServiceInstallTokenUsecase) CleanInstallToken(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return u.repository.CleanServiceInstallToken(ctx, token)
}
