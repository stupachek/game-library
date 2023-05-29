package user

import (
	"game-library/domens/models"
	"game-library/domens/repository/user_repo"
	"testing"

	"github.com/google/uuid"
)

var errs error

func BenchmarkRegister(b *testing.B) {
	repo := user_repo.NewUserRepo()
	serviсe := NewUserService(repo)
	var err error

	for i := 0; i < b.N; i++ {
		email := uuid.New().String()
		err = serviсe.Register(models.RegisterModel{
			Username: "",
			Email:    email,
			Password: "",
		})
	}

	errs = err
}
