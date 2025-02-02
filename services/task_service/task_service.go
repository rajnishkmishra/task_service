package task_service

import (
	"errors"
	"net/http"
	"time"

	"bitbucket.org/task_service/models"
	"bitbucket.org/task_service/models/vm"
	"bitbucket.org/task_service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

func (t *TaskService) ListAllTasks(ctx *utils.Context, request vm.ListAllTasksRequest) (response vm.ListAllTaskResponse, werr utils.WrapperError) {
	totalRecord := int64(0)
	dbQuery := t.db.Table("tasks")
	err := dbQuery.Count(&totalRecord).Error
	if err != nil {
		logrus.Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.TotalRecord = totalRecord
	limit := request.GetLimit()
	pages, rem := totalRecord/limit, totalRecord%limit
	totalPages := pages
	if rem != 0 {
		totalPages++
	}
	response.TotalPages = totalPages
	pageNumber := request.GetPageNumber()
	response.PageNumber = pageNumber
	offset := request.GetLimit() * (pageNumber - 1)
	dbTasks := []models.Task{}
	err = dbQuery.Offset(int(offset)).Limit(int(limit)).Find(&dbTasks).Error
	if err != nil {
		logrus.Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}

	response.Tasks = make([]vm.Task, 0)
	for _, task := range dbTasks {
		response.Tasks = append(response.Tasks, vm.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Unix(),
			UpdatedAt:   task.UpdatedAt.Unix(),
		})
	}
	return
}

func (t *TaskService) GetTaskByID(ctx *utils.Context, request vm.IDRequest) (response vm.TaskResponse, werr utils.WrapperError) {
	if request.ID <= 0 {
		err := errors.New("invalid id")
		logrus.WithContext(ctx.Ctx).Error(err)
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}
	dbTask := models.Task{}
	err := t.db.Table("tasks").Where("id = ?", request.ID).First(&dbTask).Error
	if err != nil {
		logrus.WithContext(ctx.Ctx).Error(err)
		if err == gorm.ErrRecordNotFound {
			werr = utils.NewWrapperError(http.StatusBadRequest, err)
			return
		}
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}

	response.ID = dbTask.ID
	response.Title = dbTask.Title
	response.Description = dbTask.Description
	response.Status = dbTask.Status
	response.CreatedAt = dbTask.CreatedAt.Unix()
	response.UpdatedAt = dbTask.UpdatedAt.Unix()
	return
}

func (t *TaskService) CreateTask(ctx *utils.Context, request vm.CreateTaskRequest) (response vm.TaskResponse, werr utils.WrapperError) {
	dbTask := models.Task{
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := t.db.Table("tasks").Create(&dbTask).Error
	if err != nil {
		logrus.Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.Task.Title = request.Title
	response.Task.Description = request.Description
	response.Task.Status = request.Status
	response.ID = dbTask.ID
	response.CreatedAt = dbTask.CreatedAt.Unix()
	response.UpdatedAt = dbTask.UpdatedAt.Unix()
	return
}

func (t *TaskService) UpdateTask(ctx *utils.Context, request vm.TaskRequest) (response vm.TaskResponse, werr utils.WrapperError) {
	task, werr := t.GetTaskByID(ctx, vm.IDRequest{
		ID: request.ID,
	})
	if werr != nil {
		return
	}

	dbTask := models.Task{
		ID:        request.ID,
		UpdatedAt: time.Now(),
	}
	if task.Title != request.Title {
		dbTask.Title = request.Title
	}
	if task.Description != request.Description {
		dbTask.Description = request.Description
	}
	if task.Status != request.Status {
		dbTask.Status = request.Status
	}
	err := t.db.Table("tasks").Where("id = ?", request.ID).Updates(dbTask).Error
	if err != nil {
		logrus.Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response.Task = request.Task
	response.UpdatedAt = dbTask.UpdatedAt.Unix()
	return
}

func (t *TaskService) DeleteTask(ctx *utils.Context, request vm.IDRequest) (response string, werr utils.WrapperError) {
	/*fmt.Println("khaskjckjasckjasckj")
	if request.ID <= 0 {
		fmt.Println("invalid id")
		err := errors.New("invalid id")
		werr = utils.NewWrapperError(http.StatusBadRequest, err)
		return
	}*/

	err := t.db.Unscoped().WithContext(ctx.Ctx).Table("tasks").Where("id = ?", request.ID).Delete(&models.Task{}).Error
	if err != nil {
		logrus.Error(err)
		werr = utils.NewWrapperError(http.StatusInternalServerError, err)
		return
	}
	response = "Task deleted successfully!"
	return
}
