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

func (u *UserService) Login(loginUser models.LoginModel) (string, error) {
	user, err := u.UserRepo.GetUserByEmail(loginUser.Email)
	if err != nil {
		return "", err
	}
	ok, err := argon2.VerifyEncoded([]byte(loginUser.Password), []byte(user.HashedPassword))
	if err != nil || !ok {
		return "", ErrUnauthenticated
	}
	return NewJWT(user.ID.String())
}

func newPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}
	return string(hashedPasword), nil

}
