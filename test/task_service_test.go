package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/task_service/models/vm"
	"bitbucket.org/task_service/services/backends"
	"bitbucket.org/task_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"bitbucket.org/task_service/models"
)

var (
	db     *gorm.DB
	router *gin.Engine
	token  string
)

func TestMain(m *testing.M) {
	os.Setenv("ENV", "TEST")
	gin.SetMode(gin.TestMode)

	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory database")
	}
	err = db.AutoMigrate(&models.Task{}, &models.User{})
	if err != nil {
		fmt.Println(err)
		panic("AutoMigrate failed")
	}

	router = gin.Default()
	backends.PathHandler(utils.Backends{
		GinEngine: router,
		MySQLConn: &utils.MySQLConn{
			DB: db,
		},
	})

	token = generateTestJWT()

	code := m.Run()

	os.Exit(code)
}

func generateTestJWT() string {
	request := vm.LoginRequest{
		PhoneNumber: "1234567894",
	}
	jsonData, _ := json.Marshal(request)

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/login", bytes.NewBuffer(jsonData))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type loginAPIResponse struct {
		Data vm.LoginResponse `json:"data"`
	}
	var loginAPIRes loginAPIResponse
	json.Unmarshal(resp.Body.Bytes(), &loginAPIRes)
	return loginAPIRes.Data.AccessToken
}

func TestListAllTasks(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/public/tasks", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateTask(t *testing.T) {
	task := vm.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.Pending,
	}
	jsonData, _ := json.Marshal(task)

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("token", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type createTaskAPIResponse struct {
		Data vm.TaskResponse `json:"data"`
	}
	var createTaskAPIResp createTaskAPIResponse
	json.Unmarshal(resp.Body.Bytes(), &createTaskAPIResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, task.Title, createTaskAPIResp.Data.Title)
	assert.Equal(t, task.Description, createTaskAPIResp.Data.Description)
	assert.Equal(t, task.Status, createTaskAPIResp.Data.Status)
}

func TestGetTaskByID(t *testing.T) {
	task := models.Task{
		Title:       "Task to Retrieve",
		Description: "Retrieve this task",
		Status:      models.Pending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&task)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("token", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type APIResponse struct {
		Data vm.TaskResponse `json:"data"`
	}
	var apiResponse APIResponse
	json.Unmarshal(resp.Body.Bytes(), &apiResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, task.Title, apiResponse.Data.Title)
	assert.Equal(t, task.Description, apiResponse.Data.Description)
	assert.Equal(t, task.Status, apiResponse.Data.Status)
}

func TestUpdateTask(t *testing.T) {
	task := models.Task{
		Title:       "Task to Update",
		Description: "Update this task",
		Status:      models.Pending,
	}
	db.Create(&task)

	taskUpdate := vm.TaskRequest{
		Task: vm.Task{
			ID:          task.ID,
			Title:       "Updated Title",
			Description: "Updated Description",
			Status:      models.Completed,
		},
	}
	jsonData, _ := json.Marshal(taskUpdate)

	req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/tasks/"+strconv.Itoa(int(task.ID)), bytes.NewBuffer(jsonData))
	req.Header.Set("token", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type APIResponse struct {
		Data vm.TaskResponse `json:"data"`
	}
	var apiResponse APIResponse
	json.Unmarshal(resp.Body.Bytes(), &apiResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, taskUpdate.Title, apiResponse.Data.Title)
	assert.Equal(t, taskUpdate.Description, apiResponse.Data.Description)
	assert.Equal(t, taskUpdate.Status, apiResponse.Data.Status)
}

func TestDeleteTask(t *testing.T) {
	task := models.Task{
		Title:       "Task to Delete",
		Description: "Delete this task",
		Status:      models.Pending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&task)

	req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("token", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	type APIResponse struct {
		Data string `json:"data"`
	}
	var apiResponse APIResponse
	json.Unmarshal(resp.Body.Bytes(), &apiResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Task deleted successfully!", apiResponse.Data)
}

func TestUnauthorizedAccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/tasks", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
