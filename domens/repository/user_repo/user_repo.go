package user_repo

import (
	"database/sql"
	"errors"
	"game-library/domens/models"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrDublicateUsername error = errors.New("user with the username is already exist")
	ErrDublicateEmail    error = errors.New("user with the email is already exist")
	ErrDublicateID       error = errors.New("user with the ID already exist")
	ErrUnknownUser       error = errors.New("unknown user")
	ErrUpdateFailed      error = errors.New("update failed")
	ErrDeleteFailed      error = errors.New("delete failed")
)

type TestUserRepo struct {
	Users map[uuid.UUID]*models.User
	sync.Mutex
}

func NewUserRepo() *TestUserRepo {
	return &TestUserRepo{
		Users: make(map[uuid.UUID]*models.User),
	}
}

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(DB *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		DB: DB,
	}
}

func (t *TestUserRepo) Setup() {
	t.CreateUser(models.User{
		ID:             uuid.UUID{111},
		Email:          "test",
		Username:       "test",
		Role:           "test",
		HashedPassword: "test",
	})
}

func (p *PostgresUserRepo) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS users(
        id UUID PRIMARY KEY,
    	email VARCHAR NOT NULL UNIQUE,
        userName VARCHAR NOT NULL,
		badgeColor VARCHAR,
		role VARCHAR NOT NULL,
        hashedPassword VARCHAR NOT NULL,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := p.DB.Query(query)
	return err
}

func (t *TestUserRepo) GetUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	for _, user := range t.Users {
		users = append(users, *user)
	}
	return users, nil

}

func (p *PostgresUserRepo) GetUsers() ([]models.User, error) {
	rows, err := p.DB.Query("SELECT id, email, username, role, badgeColor, createdAt FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.BadgeColor, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (t *TestUserRepo) Delete(id uuid.UUID) error {
	delete(t.Users, id)
	return nil
}

func (p *PostgresUserRepo) Delete(id uuid.UUID) error {
	res, err := p.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return nil
}

func (t *TestUserRepo) UpdateRole(id uuid.UUID, role string) error {
	user, err := t.GetUserById(id)
	if err != nil {
		return err
	}
	user.Role = role
	return nil
}

func (p *PostgresUserRepo) UpdateRole(id uuid.UUID, role string) error {
	res, err := p.DB.Exec("UPDATE users SET role = $1 WHERE  id = $2", role, id)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrUpdateFailed
	}
	return nil
}

func (t *TestUserRepo) GetUserByEmail(email string) (models.User, error) {
	for _, user := range t.Users {
		if user.Email == email {
			return *user, nil
		}
	}
	return models.User{}, ErrUnknownUser
}

func (p *PostgresUserRepo) GetUserByEmail(email string) (models.User, error) {
	row := p.DB.QueryRow("SELECT id, email, username, hashedPassword, role, badgeColor, createdAt FROM users WHERE email = $1", email)
	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.HashedPassword, &user.Role, &user.BadgeColor, &user.CreatedAt); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (t *TestUserRepo) GetUserById(id uuid.UUID) (*models.User, error) {
	user, ok := t.Users[id]
	if !ok {
		return &models.User{}, ErrUnknownUser
	}
	return user, nil
}

func (p *PostgresUserRepo) GetUserById(id uuid.UUID) (*models.User, error) {
	row := p.DB.QueryRow("SELECT id, email, username, hashedPassword,  role, badgeColor, createdAt FROM users WHERE id = $1", id)
	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.HashedPassword, &user.Role, &user.BadgeColor, &user.CreatedAt); err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (t *TestUserRepo) CreateUser(user models.User) error {
	err := t.checkUser(user)
	if err != nil {
		return err
	}
	t.Users[user.ID] = &user
	return nil
}

func (p *PostgresUserRepo) CreateUser(user models.User) error {
	_, err := p.DB.Exec("INSERT INTO users(id, email, username, badgeColor, role, hashedPassword) values($1, $2, $3,  $4, $5, $6)", user.ID, user.Email, user.Username, user.BadgeColor, user.Role, user.HashedPassword)
	return err
}

func (p *PostgresUserRepo) CreateAdmin(user models.User) error {
	_, err := p.DB.Exec("INSERT INTO users(id, email, username, badgeColor, role, hashedPassword) values($1, $2, $3,  $4, $5, $6) ON CONFLICT DO NOTHING", user.ID, user.Email, user.Username, user.BadgeColor, user.Role, user.HashedPassword)
	return err
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
