package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestNewAccountCreated(t *testing.T) {
    r := setupRouter()
    r.POST("/account", postAccount)
    account := Account{
        AccountID: "abcdefhilmn123456789",
        Name: "Nome",
        Surname: "Cognome",
        Balance: 3500,
    }
	
    jsonValue, _ := json.Marshal(account)
    req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}