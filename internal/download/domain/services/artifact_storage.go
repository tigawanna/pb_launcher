package services

import (
	"context"
	"io"

	"github.com/hashicorp/go-version"
)

type ArtifactStorage interface {
	// Versions returns a list of all available versions.
	Versions(ctx context.Context) ([]*version.Version, error)

	// Save stores the provided binary data at the specified relative path.
	Save(ctx context.Context, relativePath string, reader io.Reader) (string, error)
}
