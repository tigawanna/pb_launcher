package services

import (
	"context"
	"io"

	"github.com/hashicorp/go-version"
)

type RepositoryArtifactStorage interface {
	// Versions retrieves a list of all available versions for a given repository.
	// It returns the versions found or an error if the operation fails.
	Versions(ctx context.Context, repositoryId string) ([]*version.Version, error)

	// Save stores the provided binary data for a specific repository.
	// The data is saved at the specified relative path, returning the full storage path or an error if the operation fails.
	Save(ctx context.Context, repositoryId string, relativePath string, reader io.Reader) (string, error)
}
