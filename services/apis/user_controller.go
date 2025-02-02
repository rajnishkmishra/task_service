package apis

import (
	"bitbucket.org/task_service/services/user_service"
	"bitbucket.org/task_service/utils"
	"github.com/gin-gonic/gin"
)

const (
	Login = "/login"
)

func NewUserController(router *gin.RouterGroup, userService *user_service.UserService) {
	router.POST(Login, utils.Controller(utils.NewOptions(userService.Login).ForPost()))
}
