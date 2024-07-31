package router

import (
	"task-manager/controllers"
	"github.com/gin-gonic/gin"
)

func Router(){
	var router = gin.Default()
	gin.SetMode(gin.DebugMode)

	router.GET("/tasks", controllers.GetAllTaks)
	router.GET("/tasks/:id", controllers.GetTaskById)
	router.POST("/tasks", controllers.CreateTask)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)


	router.Run("localhost:5000")
}