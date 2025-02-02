package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/task_service/models/vm"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	task := vm.LoginRequest{
		PhoneNumber: "1010101010",
	}
	jsonData, _ := json.Marshal(task)

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type APIResponse struct {
		Data vm.LoginResponse `json:"data"`
	}
	var apiResp APIResponse
	json.Unmarshal(resp.Body.Bytes(), &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotEmpty(t, apiResp.Data.AccessToken)
}
