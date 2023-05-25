//go:build unit_test

package user

import (
	"game-library/domens/models"
	"game-library/domens/repository/user_repo"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
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
		},
		{
			description: "error dublicate",
			registerUser: models.RegisterModel{
				Username: "TEST",
				Email:    "test@test.com",
				Password: "test",
			},
			expectedError: user_repo.ErrDublicateEmail,
		},
	}
	repo := user_repo.NewUserRepo()
	servive := NewUserService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := servive.Register(tc.registerUser)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}

func TestLogin(t *testing.T) {
	testCases := []struct {
		description   string
		loginUser     models.LoginModel
		expectedError error
	}{
		{
			description: "wrong password",
			loginUser: models.LoginModel{
				Email:    "test",
				Password: "test2",
			},
			expectedError: ErrUnauthenticated,
		},
		{
			description: "unknown user",
			loginUser: models.LoginModel{
				Email:    "test2",
				Password: "test2",
			},
			expectedError: ErrUnauthenticated,
		},
	}
	repo := user_repo.NewUserRepo()
	err := repo.Setup()
	if err != nil {
		t.Error(err)
	}
	servive := NewUserService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := servive.Login(tc.loginUser)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		expected      models.User
		expectedError error
	}{
		{
			description: "success",
			idStr:       uuid.UUID{111}.String(),
			expected: models.User{
				ID:             uuid.UUID{111},
				Email:          "test",
				Username:       "test",
				Role:           "test",
				HashedPassword: "test",
			},
			expectedError: nil,
		},
		{
			description:   "can't parse uuid",
			idStr:         "smth",
			expected:      models.User{},
			expectedError: ErrUnknownUUID,
		},
	}
	repo := user_repo.NewUserRepo()
	err := repo.Setup()
	if err != nil {
		t.Error(err)
	}
	servive := NewUserService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			user, err := servive.GetUser(tc.idStr)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if !reflect.DeepEqual(tc.expected, user) {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expected, user)
			}
		})

	}

}

func TestGetDeleteUser(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		expectedError error
	}{
		{
			description:   "success",
			idStr:         uuid.UUID{111}.String(),
			expectedError: nil,
		},
		{
			description:   "can't parse uuid",
			idStr:         "smth",
			expectedError: ErrUnknownUUID,
		},
	}
	repo := user_repo.NewUserRepo()
	err := repo.Setup()
	if err != nil {
		t.Error(err)
	}
	servive := NewUserService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := servive.DeleteUser(tc.idStr)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
		})

	}

}

func TestChangeRole(t *testing.T) {
	testCases := []struct {
		description   string
		idStr         string
		role          string
		expected      models.User
		expectedError error
	}{
		{
			description: "success admin",
			idStr:       uuid.UUID{111}.String(),
			role:        "admin",
			expected: models.User{
				ID:             uuid.UUID{111},
				Email:          "test",
				Username:       "test",
				Role:           "admin",
				HashedPassword: "test",
			},
			expectedError: nil,
		},
		{
			description: "success user",
			idStr:       uuid.UUID{111}.String(),
			role:        "user",
			expected: models.User{
				ID:             uuid.UUID{111},
				Email:          "test",
				Username:       "test",
				Role:           "user",
				HashedPassword: "test",
			},
			expectedError: nil,
		},
		{
			description: "success manager",
			idStr:       uuid.UUID{111}.String(),
			role:        "manager",
			expected: models.User{
				ID:             uuid.UUID{111},
				Email:          "test",
				Username:       "test",
				Role:           "manager",
				HashedPassword: "test",
			},
			expectedError: nil,
		},
		{
			description:   "unknown role",
			idStr:         uuid.UUID{111}.String(),
			role:          "role",
			expectedError: ErrUnknownRole,
		},
		{
			description:   "can't parse uuid",
			idStr:         "smth",
			expectedError: ErrUnknownUUID,
		},
	}
	repo := user_repo.NewUserRepo()
	err := repo.Setup()
	if err != nil {
		t.Error(err)
	}
	servive := NewUserService(repo)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			user, err := servive.ChangeRole(tc.idStr, tc.role)
			if err != tc.expectedError {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedError, err)
			}
			if !reflect.DeepEqual(tc.expected, user) {
				t.Fatalf("%s. expected %v, got %v", tc.description, tc.expected, user)
			}
		})

	}

}
