package services

import (
	"context"
	"regexp"
)

type BinaryFinder interface {
	FindBinary(ctx context.Context, repositoryID, version string, binaryPattern *regexp.Regexp) (string, error)
}
