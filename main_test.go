package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountCreated_Integration_1(t *testing.T) {
    r := setupRouter()
    r.POST("/account", postAccount)
    account := Account{
        Name: "NomeProvaMittente",
        Surname: "NomeProvaMittente",
        Balance: 10,
    }
	
    jsonValue, _ := json.Marshal(account)
    req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPingRoute_Unit_1(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
