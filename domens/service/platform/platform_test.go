//go:build unit_test

package platform

import (
	"game-library/domens/models"
	"game-library/domens/repository/platform_repo"
	"testing"

	"github.com/google/uuid"
)

func TestWCreateSucces(t *testing.T) {
	repo := platform_repo.NewPlatformRepo()
	service := NewPlatformService(repo)
	t.Run("success create platform", func(t *testing.T) {
		platform, err := service.CreatePlatform(models.Platform{
			Name: "test",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if platform.ID == (uuid.UUID{}) {
			t.Fatal("got empty uuid")
		}
	})
}

func TestGetSucces(t *testing.T) {
	testCases := []struct {
		description          string
		platform             string
		expectedError        error
		expectedplatformName string
	}{
		{
			description:          "success",
			platform:             "test",
			expectedError:        nil,
			expectedplatformName: "test",
		},
		{
			description:          "failed",
			platform:             "error",
			expectedError:        ErrUnknownPlatform,
			expectedplatformName: "",
		},
	}
	repo := platform_repo.NewPlatformRepo()
	repo.Setup()
	service := NewPlatformService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			platform, err := service.GetPlatform(tc.platform)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if platform.Name != tc.expectedplatformName {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedplatformName, platform.Name)
			}
		})

	}
}

func TestGetListplatforms(t *testing.T) {
	repo := platform_repo.NewPlatformRepo()
	repo.Setup()
	service := NewPlatformService(repo)
	t.Run("success get platforms", func(t *testing.T) {
		platforms, err := service.GetPlatformsList()
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if len(platforms) != 1 {
			t.Fatalf("expected %v, got %v", 1, len(platforms))
		}
	})

}
