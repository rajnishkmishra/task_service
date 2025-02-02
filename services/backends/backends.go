package backends

import (
	"os"

	"bitbucket.org/task_service/migration"
	"bitbucket.org/task_service/services/apis"
	"bitbucket.org/task_service/services/auth_service"
	"bitbucket.org/task_service/services/task_service"
	"bitbucket.org/task_service/services/user_service"
	"bitbucket.org/task_service/utils"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func PathHandler(backends utils.Backends) {
	if os.Getenv("ENV") != "TEST" {
		migrationService := migration.NewMigrationService(backends.MySQLConn)
		go migrationService.InitMigration()
	}
	r = backends.GinEngine
	v1 := r.Group("/v1")

	userService := user_service.NewUserService(backends.MySQLConn.DB)
	authService := auth_service.NewAuthService(userService)

	apis.NewUserController(v1, userService)
	apis.NewTaskController(v1, task_service.NewTaskService(backends.MySQLConn.DB), authService)
}
