package test

import (
	"app/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func registerUser(t *testing.T, name, email string, password string) int {
	registerRequest := &domain.RegisterRequestDTO{
		Name:     name,
		Email:    email,
		Password: password,
	}
	jsonValue, err := json.Marshal(registerRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.NotZero(t, response.Data["id"])
	return int(response.Data["id"].(float64))
}

func loginUser(t *testing.T, email string, password string) string {
	loginRequest := &domain.LoginRequestDTO{
		Email:    email,
		Password: password,
	}
	jsonValue, err := json.Marshal(loginRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)
	var cookieValue string
	for _, cookie := range cookies {
		if cookie.Name == "AUTHORIZATION" {
			cookieValue = cookie.Value
		}
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	return cookieValue
}

func TestLogin(t *testing.T) {
	t.Run("login with unregistered email", func(t *testing.T) {
		loginRequest := &domain.LoginRequestDTO{
			Email:    "unregistered@email.com",
			Password: "password",
		}
		jsonValue, err := json.Marshal(loginRequest)
		assert.Nil(t, err)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "Email not found", response["message"])
	})
	t.Run("login with wrong password", func(t *testing.T) {
		registerUser(t, "name", "email@email.com", "password")
		loginRequest := &domain.LoginRequestDTO{
			Email:    "email@email.com",
			Password: "wrongpassword",
		}
		jsonValue, err := json.Marshal(loginRequest)
		assert.Nil(t, err)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid password", response["message"])
	})

	t.Run("login with correct email and password", func(t *testing.T) {
		id := registerUser(t, "name", "email1@email.com", "password")
		cookie := loginUser(t, "email1@email.com", "password")
		assert.NotEmpty(t, cookie)
		assert.NotZero(t, id)
	})
}
