//go:build api_test

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game-library/domens/models"
	"game-library/domens/repository/database"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
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

func TestGenres(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(DB)

	t.Run("genres", func(t *testing.T) {
		//login admin
		inputL := models.LoginModel{
			Email:    "admin@a.a",
			Password: "admin",
		}
		jsonValue, _ := json.Marshal(inputL)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())
		token := w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create genre
		input := models.Genre{
			Name: "test",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/genres", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"test\"},\"message\":\"Genre is successfully created\"}", w.Body.String())

		// sing up user
		inputR := models.RegisterModel{
			Username: "test",
			Email:    "test@test.com",
			Password: "test",
		}
		jsonValue, _ = json.Marshal(inputR)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\"message\":\"Sign up was successful\"}", w.Body.String())

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
		token = w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create genre without permissions
		input = models.Genre{
			Name: "testUser",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/genres", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "{\"error\":\"permission denied\"}", w.Body.String())

		//get genres
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/genres", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "\"name\":\"test\"}]}$", w.Body.String())

	})
}

func TestPublishers(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(DB)

	t.Run("publisher", func(t *testing.T) {
		//login admin
		inputL := models.LoginModel{
			Email:    "admin@a.a",
			Password: "admin",
		}
		jsonValue, _ := json.Marshal(inputL)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())
		token := w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create publisher
		input := models.Publisher{
			Name: "test",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/publishers", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"test\"},\"message\":\"Publisher is successfully created\"}", w.Body.String())

		// sing up user
		inputR := models.RegisterModel{
			Username: "test",
			Email:    "test@test.com",
			Password: "test",
		}
		jsonValue, _ = json.Marshal(inputR)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\"message\":\"Sign up was successful\"}", w.Body.String())

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
		token = w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create publisher without permissions
		input = models.Publisher{
			Name: "testUser",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/publishers", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "{\"error\":\"permission denied\"}", w.Body.String())

		//get publishers
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/publishers", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "\"name\":\"test\"}]}$", w.Body.String())

	})
}

func TestPlatforms(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(DB)

	t.Run("platforms", func(t *testing.T) {
		//login admin
		inputL := models.LoginModel{
			Email:    "admin@a.a",
			Password: "admin",
		}
		jsonValue, _ := json.Marshal(inputL)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())
		token := w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create platform
		input := models.Platform{
			Name: "test",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/platforms", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"test\"},\"message\":\"Platform is successfully created\"}", w.Body.String())

		// sing up user
		inputR := models.RegisterModel{
			Username: "test",
			Email:    "test@test.com",
			Password: "test",
		}
		jsonValue, _ = json.Marshal(inputR)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\"message\":\"Sign up was successful\"}", w.Body.String())

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
		token = w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create platform without permissions
		input = models.Platform{
			Name: "test2",
		}
		jsonValue, _ = json.Marshal(input)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/platforms", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "{\"error\":\"permission denied\"}", w.Body.String())

		//get platforms
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/platforms", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "\"name\":\"test\"}]}$", w.Body.String())

	})
}

func TestGamesSuccess(t *testing.T) {
	DB := database.ConnectDataBase()
	err := database.ClearData(DB)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(DB)

	t.Run("create, get game", func(t *testing.T) {
		//login admin
		inputL := models.LoginModel{
			Email:    "admin@a.a",
			Password: "admin",
		}
		jsonValue, _ := json.Marshal(inputL)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonValue))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"message\":\"Sign up was successful\",\"token\":\"([a-zA-Z0-9-_.]{207})\"}", w.Body.String())
		token := w.Body.String()[w.Body.Len()-209 : w.Body.Len()-2]

		//create platform
		inputPlatform := models.Platform{
			Name: "testPlatform",
		}
		jsonValue, _ = json.Marshal(inputPlatform)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/platforms", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"testPlatform\"},\"message\":\"Platform is successfully created\"}", w.Body.String())

		//create genre
		inputGenre := models.Genre{
			Name: "testGenre1",
		}
		jsonValue, _ = json.Marshal(inputGenre)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/genres", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"testGenre1\"},\"message\":\"Genre is successfully created\"}", w.Body.String())

		//create genre
		inputGenre = models.Genre{
			Name: "testGenre2",
		}
		jsonValue, _ = json.Marshal(inputGenre)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/genres", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"testGenre2\"},\"message\":\"Genre is successfully created\"}", w.Body.String())

		//create publisher
		inputPublisher := models.Publisher{
			Name: "test",
		}
		jsonValue, _ = json.Marshal(inputPublisher)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/publishers", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"id\":\"([a-zA-Z0-9-]{36})\",\"name\":\"test\"},\"message\":\"Publisher is successfully created\"}", w.Body.String())
		pattern := regexp.MustCompile("[a-zA-Z0-9-]{36}")
		publisherId := pattern.FindString(w.Body.String())

		genres := []string{"testGenre1", "testGenre2"}
		platforms := []string{"testPlatform"}
		fmt.Print(genres, platforms)

		//create game
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fw, err := writer.CreateFormField("title")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, strings.NewReader("testGame"))
		if err != nil {
			t.Error(err)
		}
		fw, err = writer.CreateFormField("publisherId")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, strings.NewReader(publisherId))
		if err != nil {
			t.Error(err)
		}

		fw, err = writer.CreateFormFile("file", "test.png")
		if err != nil {
		}
		file, err := os.Open("test.png")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			t.Error(err)
		}

		fw, err = writer.CreateFormField("genres")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, strings.NewReader(genres[0]))
		if err != nil {
			t.Error(err)
		}
		fw, err = writer.CreateFormField("genres")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, strings.NewReader(genres[1]))
		if err != nil {
			t.Error(err)
		}
		fw, err = writer.CreateFormField("platforms")
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(fw, strings.NewReader(platforms[0]))
		if err != nil {
			t.Error(err)
		}

		writer.Close()

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/games", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Regexp(t, "{\"data\":{\"gameId\":\"([a-zA-Z0-9-]{36})\",\"link\":\"http://localhost:8080/image/library/test.png\"},\"message\":\"Game is successfully created\"", w.Body.String())
		err = os.Remove("library/test.png")
		if err != nil {
			t.Error(err)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/games", bytes.NewReader(body.Bytes()))
		req, _ = http.NewRequest("GET", "/games", bytes.NewReader(body.Bytes()))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

	})
}
