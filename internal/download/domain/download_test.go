package domain_test

import (
	"pb_luncher/internal/download/domain"
	"pb_luncher/internal/download/domain/dtos"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestDiffReleases(t *testing.T) {
	uc := &domain.DownloadUsecase{}

	v1, _ := version.NewVersion("1.0.0")
	v2, _ := version.NewVersion("1.2.0")
	v3, _ := version.NewVersion("2.0.0")
	v4, _ := version.NewVersion("3.0.0")

	a := []dtos.Release{
		{Version: v1},
		{Version: v2},
		{Version: v3},
		{Version: v4},
	}

	b := []dtos.Release{
		{Version: v1},
		{Version: v3},
	}

	expected := []dtos.Release{
		{Version: v2},
		{Version: v4},
	}

	diff := uc.DiffReleases(a, b)

	if len(diff) != len(expected) {
		t.Errorf("Expected %d releases, got %d", len(expected), len(diff))
	}

	for i, r := range diff {
		if !r.Version.Equal(expected[i].Version) {
			t.Errorf("Expected version %s, got %s", expected[i].Version, r.Version)
		}
	}
}
