//go:build integration_test

package user

import (
	"crypto/ed25519"
	"errors"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"game-library/domens/repository/user_repo"
	"game-library/domens/service/jwt"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestRegisterLogin(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := user_repo.NewPostgresUserRepo(DB)
	userService := NewUserService(userRepo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	testCases := []struct {
		description    string
		registerUser   models.RegisterModel
		loginUser      models.LoginModel
		expectedErrorR error
		expectedErrorL error
	}{
		{
			description: "success register, login",
			registerUser: models.RegisterModel{
				Username: "test",
				Email:    "test",
				Password: "test",
			},
			loginUser: models.LoginModel{
				Email:    "test",
				Password: "test",
			},
			expectedErrorR: nil,
			expectedErrorL: nil,
		},
		{
			description:    "dublicate",
			registerUser:   models.RegisterModel{Username: "test", Email: "test", Password: "test2"},
			loginUser:      models.LoginModel{Email: "test", Password: "test2"},
			expectedErrorR: errors.New("duplicate key value violates"),
			expectedErrorL: ErrUnauthenticated,
		},
		{
			description: "wrong password",
			registerUser: models.RegisterModel{
				Username: "test1",
				Email:    "test1",
				Password: "test1",
			},
			loginUser: models.LoginModel{
				Email:    "test1",
				Password: "test2",
			},
			expectedErrorR: nil,
			expectedErrorL: ErrUnauthenticated,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := userService.Register(tc.registerUser)
			if err != tc.expectedErrorR {
				if tc.expectedErrorR == nil {
					t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedErrorR, err)
				} else if !strings.Contains(err.Error(), tc.expectedErrorR.Error()) {
					t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedErrorR, err)
				}
			}
			_, err = userService.Login(tc.loginUser)
			if err != tc.expectedErrorL {
				if tc.expectedErrorL == nil {
					t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedErrorL, err)
				} else if !strings.Contains(err.Error(), tc.expectedErrorL.Error()) {
					t.Fatalf("%s. expected %v, got %v", tc.description, tc.expectedErrorL, err)
				}
			}
		})

	}

}

func TestGetChangeRoleUser(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := user_repo.NewPostgresUserRepo(DB)
	userService := NewUserService(userRepo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("get, change, delete user", func(t *testing.T) {
		err := userService.Register(models.RegisterModel{
			Username: "test",
			Email:    "test",
			Password: "test",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		err = userService.Register(models.RegisterModel{
			Username: "test2",
			Email:    "test2",
			Password: "test2",
		})
		err = userService.Register(models.RegisterModel{
			Username: "test3",
			Email:    "test3",
			Password: "test3",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		users, err := userService.GetUsers()
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if len(users) != 4 {
			t.Fatalf(" expected %v users, got %v", 4, len(users))
		}
		user, err := userService.ChangeRole(users[0].ID.String(), models.ADMIN)
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		user, err = userService.GetUser(users[0].ID.String())
		if user.Role != models.ADMIN {
			t.Fatalf(" expected %v, got %v", models.ADMIN, user.Role)
		}
		user, err = userService.ChangeRole(users[1].ID.String(), models.MANAGER)
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		user, err = userService.GetUser(users[1].ID.String())
		if user.Role != models.MANAGER {
			t.Fatalf(" expected %v, got %v", models.MANAGER, user.Role)
		}
		user, err = userService.ChangeRole(users[2].ID.String(), "SMTH")
		if err != ErrUnknownRole {
			t.Fatalf(" expected %v, got %v", ErrUnknownRole, err)
		}
		user, err = userService.GetUser(users[2].ID.String())
		if user.Role != models.USER {
			t.Fatalf(" expected %v, got %v", models.USER, user.Role)
		}
	})

}

func TestDeleteUser(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := user_repo.NewPostgresUserRepo(DB)
	userService := NewUserService(userRepo)
	jwt.Public, jwt.Private, err = ed25519.GenerateKey(nil)
	t.Run("get, change, delete user", func(t *testing.T) {
		err := userService.Register(models.RegisterModel{
			Username: "test",
			Email:    "test",
			Password: "test",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		err = userService.Register(models.RegisterModel{
			Username: "test2",
			Email:    "test2",
			Password: "test2",
		})
		err = userService.Register(models.RegisterModel{
			Username: "test3",
			Email:    "test3",
			Password: "test3",
		})
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		users, err := userService.GetUsers()
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if len(users) != 4 {
			t.Fatalf(" expected %v users, got %v", 4, len(users))
		}
		err = userService.DeleteUser(users[0].ID.String())
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		users, err = userService.GetUsers()
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if len(users) != 3 {
			t.Fatalf(" expected %v users, got %v", 3, len(users))
		}
		err = userService.DeleteUser("error")
		if err != ErrUnknownUUID {
			t.Fatalf(" expected %v, got %v", ErrUnknownUUID, err)
		}
		users, err = userService.GetUsers()
		if err != nil {
			t.Fatalf(" expected %v, got %v", nil, err)
		}
		if len(users) != 3 {
			t.Fatalf(" expected %v users, got %v", 3, len(users))
		}

	})

}
