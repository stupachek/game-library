package repository

import (
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	dublicateUsernameError error = errors.New("user with the username is already exist")
	dublicateEmailError    error = errors.New("user with the email is already exist")
	dublicateIDError       error = errors.New("user with the ID is already exist")
)

type IUserRepo interface {
	CreateUser(models.User) error
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
		return dublicateIDError
	}
	for _, u := range t.Users {
		if u.Email == user.Email {
			return dublicateEmailError
		} else if u.Username == user.Username {
			return dublicateUsernameError
		}
	}
	return nil
}
