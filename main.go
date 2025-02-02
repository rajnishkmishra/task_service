package main

import (
	"bitbucket.org/task_service/services/backends"
	"bitbucket.org/task_service/utils"
)

func main() {
	utils.SetupAndRun(backends.PathHandler)
}
