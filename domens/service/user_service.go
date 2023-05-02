package service

import (
	"errors"
	"game-library/domens/models"
	"game-library/domens/repository"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrAdmin           = errors.New("setup admin error")
)

type UserService struct {
	UserRepo repository.IUserRepo
}

func NewUserService(repo repository.IUserRepo) UserService {
	return UserService{
		UserRepo: repo,
	}
}

func (u *UserService) SetupAdmin() error {
	ADMIN_EMAIL, ok := os.LookupEnv("ADMIN_EMAIL")
	if !ok {
		log.Fatal("please specify ADMIN_EMAIL")
	}
	ADMIN_USERNAME, ok := os.LookupEnv("ADMIN_USERNAME")
	if !ok {
		log.Fatal("please specify ADMIN_USERNAME")
	}
	ADMIN_PASSWORD, ok := os.LookupEnv("ADMIN_PASSWORD")
	if !ok {
		log.Fatal("please specify ADMIN_PASSWORD")
	}
	hashedPassword, err := newPassword(ADMIN_PASSWORD)
	if err != nil {
		return ErrAdmin
	}
	user := models.User{
		Email:          ADMIN_EMAIL,
		Username:       ADMIN_USERNAME,
		HashedPassword: hashedPassword,
		Role:           models.ADMIN,
	}
	user.ID, err = uuid.NewRandom()
	if err != nil {
		return ErrAdmin
	}
	err = u.UserRepo.CreateAdmin(user)
	return err
}

func (u *UserService) Register(registerUser models.RegisterModel) error {
	hashedPassword, err := newPassword(registerUser.Password)
	if err != nil {
		return ErrUnauthenticated
	}
	user := models.User{
		Email:          registerUser.Email,
		Username:       registerUser.Username,
		Role:           models.USER,
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

func (u *UserService) GetUser(idStr string) (models.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.User{}, errors.New("can't parse user id")
	}
	user, err := u.UserRepo.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}
	return *user, err
}

func (u *UserService) DeleteUser(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("can't parse user id")
	}
	if _, err := u.UserRepo.GetUserById(id); err != nil {
		return err
	}
	return u.UserRepo.Delete(id)
}

func (u *UserService) GetUsers() ([]models.User, error) {
	users, err := u.UserRepo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) ChangeRole(idStr string, role string) (models.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.User{}, errors.New("can't parse user id")
	}
	var user models.User
	switch role {
	case models.USER, models.ADMIN, models.MANAGER:
		if err := u.UserRepo.UpdateRole(id, role); err != nil {
			return models.User{}, err
		}
		userP, err := u.UserRepo.GetUserById(id)
		if err != nil {
			return models.User{}, err
		}
		user = *userP
	default:
		return models.User{}, errors.New("unknown role")
	}
	return user, nil
}

func newPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}
	return string(hashedPasword), nil

}
