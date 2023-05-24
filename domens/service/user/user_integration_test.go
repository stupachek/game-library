//go:build integration_test

package user

import (
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/user_repo"
	"testing"
)

func TestI(t *testing.T) {
	DB := database.ConnectDataBase()
	userRepo := user_repo.NewPostgresUserRepo(DB)
	userService := NewUserService(userRepo)
	testCases := []struct {
		description   string
		registerUser  models.RegisterModel
		expectedError error
	}{
		{
			description: "success",
			registerUser: models.RegisterModel{
				Username: "TEST",
				Email:    "test@test.com",
				Password: "test",
			},
			expectedError: nil,
		}}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := userService.Register(tc.registerUser)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}
