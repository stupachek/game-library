package service

import (
	"errors"
	"game-library/domens/models"
	"game-library/domens/repository"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

const (
	USER = "user"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

type UserService struct {
	UserRepo repository.IUserRepo
}

func NewUserService(repo repository.IUserRepo) UserService {
	return UserService{
		UserRepo: repo,
	}
}

func (u *UserService) Register(registerUser models.RegisterModel) error {
	hashedPassword, err := newPassword(registerUser.Password)
	if err != nil {
		return ErrUnauthenticated
	}
	user := models.User{
		Email:          registerUser.Email,
		Username:       registerUser.Username,
		BadgeColor:     hashedPassword,
		Role:           USER,
		HashedPassword: hashedPassword,
	}
	user.ID, err = uuid.NewRandom()
	if err != nil {
		return ErrUnauthenticated
	}
	err = u.UserRepo.CreateUser(user)
	return err
}

func newPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}
	return string(hashedPasword), nil

}
