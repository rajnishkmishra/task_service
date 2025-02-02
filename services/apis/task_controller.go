package apis

import (
	"bitbucket.org/task_service/services/auth_service"
	"bitbucket.org/task_service/services/task_service"
	"bitbucket.org/task_service/utils"
	"github.com/gin-gonic/gin"
)

const (
	ListAllTasks = "/public/tasks"
	GetTaskByID  = "/tasks/:id"
	CreateTask   = "/tasks"
	UpdateTask   = "/tasks/:id"
	DeleteTask   = "/tasks/:id"
)

func NewTaskController(router *gin.RouterGroup, taskService *task_service.TaskService, authService *auth_service.AuthService) {
	router.GET(ListAllTasks, utils.Controller(utils.NewOptions(taskService.ListAllTasks)))
	router.Use(authService.AuthMiddleWare())
	router.Use(authService.RateLimit())
	router.GET(GetTaskByID, utils.Controller(utils.NewOptions(taskService.GetTaskByID)))
	router.POST(CreateTask, utils.Controller(utils.NewOptions(taskService.CreateTask).ForPost()))
	router.PUT(UpdateTask, utils.Controller(utils.NewOptions(taskService.UpdateTask).ForPost()))
	router.DELETE(DeleteTask, utils.Controller(utils.NewOptions(taskService.DeleteTask)))
}
