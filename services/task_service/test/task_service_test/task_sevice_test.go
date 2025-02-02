package task_service_test

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/task_service/models"
	"bitbucket.org/task_service/models/vm"
	"bitbucket.org/task_service/services/task_service"
	"bitbucket.org/task_service/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&models.Task{})
	os.Exit(m.Run())
}

func TestListAllTasks(t *testing.T) {
	taskService := task_service.NewTaskService(db)

	db.Create(&models.Task{Title: "Task 1", Description: "Description 1", Status: models.Pending})
	db.Create(&models.Task{Title: "Task 2", Description: "Description 2", Status: models.Completed})

	request := vm.ListAllTasksRequest{
		PaginationRequest: vm.NewPaginationRequest(1, 10),
	}

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	response, werr := taskService.ListAllTasks(ctx, request)

	assert.Nil(t, werr)
	assert.Equal(t, int64(2), response.TotalRecord)
	assert.Equal(t, int64(1), response.PageNumber)
	assert.Equal(t, int64(1), response.TotalPages)
	assert.Len(t, response.Tasks, 2)
	assert.Equal(t, "Task 1", response.Tasks[0].Title)
	assert.Equal(t, "Task 2", response.Tasks[1].Title)
}

func TestCreateTask(t *testing.T) {
	taskService := task_service.NewTaskService(db)

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	request := vm.CreateTaskRequest{
		Title:       "New Task",
		Description: "New task description",
		Status:      models.Pending,
	}

	response, werr := taskService.CreateTask(ctx, request)

	assert.Nil(t, werr)

	assert.Equal(t, "New Task", response.Title)
	assert.Equal(t, "New task description", response.Description)
	assert.Equal(t, models.Pending, response.Status)

	var task models.Task
	err := db.Table("tasks").Where("id = ?", response.ID).First(&task).Error
	assert.Nil(t, err)
	assert.Equal(t, "New Task", task.Title)
	assert.Equal(t, "New task description", task.Description)
	assert.Equal(t, models.Pending, task.Status)
}

func TestGetTaskByID(t *testing.T) {
	taskService := task_service.NewTaskService(db)

	task := models.Task{Title: "Task 1", Description: "Description 1", Status: models.Pending}
	db.Create(&task)

	request := vm.IDRequest{
		ID: task.ID,
	}

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	response, werr := taskService.GetTaskByID(ctx, request)

	assert.Nil(t, werr)
	assert.Equal(t, task.ID, response.ID)
	assert.Equal(t, "Task 1", response.Title)
	assert.Equal(t, "Description 1", response.Description)
	assert.Equal(t, models.Pending, response.Status)
}

func TestUpdateTask(t *testing.T) {
	taskService := task_service.NewTaskService(db)

	db.Create(&models.Task{Title: "Task 1", Description: "Description 1", Status: models.Pending})

	request := vm.TaskRequest{
		Task: vm.Task{
			ID:          1,
			Title:       "Updated Task 1",
			Description: "Updated Description 1",
			Status:      models.Completed,
		},
	}

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	response, werr := taskService.UpdateTask(ctx, request)

	assert.Nil(t, werr)
	assert.Equal(t, "Updated Task 1", response.Title)
	assert.Equal(t, "Updated Description 1", response.Description)
	assert.Equal(t, models.Completed, response.Status)

	var task models.Task
	err := db.Table("tasks").Where("id = ?", response.ID).First(&task).Error
	assert.Nil(t, err)
	assert.Equal(t, "Updated Task 1", task.Title)
	assert.Equal(t, "Updated Description 1", task.Description)
	assert.Equal(t, models.Completed, task.Status)
}

func TestDeleteTask(t *testing.T) {
	taskService := task_service.NewTaskService(db)

	db.Create(&models.Task{Title: "Task 1", Description: "Description 1", Status: models.Pending})

	request := vm.IDRequest{
		ID: 1,
	}

	ctx, cancel := utils.CreateBackgroundContextWithTimeout(10 * time.Second)
	defer cancel()
	response, werr := taskService.DeleteTask(ctx, request)

	assert.Nil(t, werr)

	assert.Equal(t, "Task deleted successfully!", response)

	var task models.Task
	err := db.Table("tasks").Where("id = ?", request.ID).First(&task).Error
	assert.NotNil(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
