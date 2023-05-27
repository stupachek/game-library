//go:build api_test

package api

import (
	"bytes"
	"encoding/json"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(DB)

	t.Run("user", func(t *testing.T) {
		// sing up user
		inputR := models.RegisterModel{
			Username: "test",
			Email:    "test@test.com",
			Password: "test",
		}
		jsonValue, _ := json.Marshal(inputR)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\"message\":\"Sign up was successful\"}", w.Body.String())

		// sing up dublicate user
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "{\"error\":\"pq: duplicate key value violates unique constraint \\\"users_email_key\\\"\",\"message\":\"can't register\"}", w.Body.String())

		// sing up can't bind
		wrong := ""
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte(wrong)))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "{\"error\":\"EOF\",\"message\":\"can't parse input\"}", w.Body.String())

		// sing in can't bind
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signin", bytes.NewBuffer([]byte(wrong)))
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		assert.Equal(t, "{\"error\":\"EOF\",\"message\":\"can't parse input\"}", w.Body.String())

		//login user with wrong password
		inputL := models.LoginModel{
			Email:    "test@test.com",
			Password: "test123",
		}
		jsonValue, _ = json.Marshal(inputL)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		assert.Equal(t, "{\"error\":\"unauthenticated\",\"message\":\"can't login\"}", w.Body.String())

		//login user
		inputL = models.LoginModel{
			Email:    "test@test.com",
			Password: "test",
		}
		jsonValue, _ = json.Marshal(inputL)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())

		//Get info about user
		token := w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/users/me", nil)
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"email\":\"test@test.com\",\"usarname\":\"test\",\"badge_color\":\"\",\"role\":\"user\"}}", w.Body.String())
		router.ServeHTTP(w, req)
		pattern := regexp.MustCompile("[a-zA-Z0-9-]{36}")
		id := pattern.FindString(w.Body.String())

		//login admin
		inputL = models.LoginModel{
			Email:    "admin@a.a",
			Password: "admin",
		}
		jsonValue, _ = json.Marshal(inputL)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())
		token = w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//change role
		inputRole := models.Role{
			Role: "manager",
		}
		jsonValue, _ = json.Marshal(inputRole)
		w = httptest.NewRecorder()
		url := "/users/" + id
		req, _ = http.NewRequest("PATCH", url, bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"email\":\"test@test.com\",\"usarname\":\"test\",\"badge_color\":\"\",\"role\":\"manager\"},\"message\":\"User is successfully updated\"}", w.Body.String())

		//delete user
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\"message\":\"User is successfully deleted\"}", w.Body.String())

		//get deleted user
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", url, nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "{\"error\":\"sql: no rows in result set\",\"message\":\"can't get user\"}", w.Body.String())

	})
}
