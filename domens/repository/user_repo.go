package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrDublicateUsername error = errors.New("user with the username is already exist")
	ErrDublicateEmail    error = errors.New("user with the email is already exist")
	ErrDublicateID       error = errors.New("user with the ID is already exist")
	ErrUnknownUser       error = errors.New("unknown user")
)

type IUserRepo interface {
	CreateUser(models.User) error
	GetUserByEmail(email string) (models.User, error)
}

type TestUserRepo struct {
	Users map[uuid.UUID]models.User
	sync.Mutex
}

func NewUserRepo() *TestUserRepo {
	return &TestUserRepo{
		Users: make(map[uuid.UUID]models.User),
	}
}

func (t *TestUserRepo) GetUserByEmail(email string) (models.User, error) {
	for _, user := range t.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, ErrUnknownUser
}

func (t *TestUserRepo) CreateUser(user models.User) error {
	err := t.checkUser(user)
	if err != nil {
		return err
	}
	t.Users[user.ID] = user
	return nil
}

func (t *TestUserRepo) checkUser(user models.User) error {
	_, ok := t.Users[user.ID]
	if ok {
		return ErrDublicateID
	}
	for _, u := range t.Users {
		if u.Email == user.Email {
			return ErrDublicateEmail
		} else if u.Username == user.Username {
			return ErrDublicateUsername
		}
	}
	return nil
}
