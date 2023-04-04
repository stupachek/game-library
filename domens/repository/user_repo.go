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
	GetUserById(id uuid.UUID) (*models.User, error)
	GetUsers() []models.User
	UpdateRole(id uuid.UUID, role string) (models.User, error)
	Delete(id uuid.UUID)
}

type TestUserRepo struct {
	Users map[uuid.UUID]*models.User
	sync.Mutex
}

func NewUserRepo() *TestUserRepo {
	return &TestUserRepo{
		Users: make(map[uuid.UUID]*models.User),
	}
}

func (t *TestUserRepo) GetUsers() []models.User {
	users := make([]models.User, 0)
	for _, user := range t.Users {
		users = append(users, *user)
	}
	return users

}

func (t *TestUserRepo) Delete(id uuid.UUID) {
	delete(t.Users, id)
}

func (t *TestUserRepo) UpdateRole(id uuid.UUID, role string) (models.User, error) {
	user, err := t.GetUserById(id)
	if err != nil {
		return *user, err
	}
	user.Role = role
	return *user, nil
}

func (t *TestUserRepo) GetUserByEmail(email string) (models.User, error) {
	for _, user := range t.Users {
		if user.Email == email {
			return *user, nil
		}
	}
	return models.User{}, ErrUnknownUser
}

func (t *TestUserRepo) GetUserById(id uuid.UUID) (*models.User, error) {
	user, ok := t.Users[id]
	if !ok {
		return &models.User{}, ErrUnknownUser
	}
	return user, nil
}

func (t *TestUserRepo) CreateUser(user models.User) error {
	err := t.checkUser(user)
	if err != nil {
		return err
	}
	t.Users[user.ID] = &user
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
