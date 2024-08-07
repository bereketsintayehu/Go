package router

import (
	"task-manager/controllers"
	"task-manager/middlewares"

	"github.com/gin-gonic/gin"
)

func Router(){
	var router = gin.Default()
	gin.SetMode(gin.DebugMode)

	router.GET("/tasks", middlewares.AuthMiddleware(),controllers.GetAllTaks)
	router.GET("/tasks/:id", middlewares.AuthMiddleware(),controllers.GetTaskById)
	router.POST("/tasks", middlewares.AuthMiddleware(),controllers.CreateTask)
	router.PUT("/tasks/:id", middlewares.AuthMiddleware(),controllers.UpdateTask)
	router.DELETE("/tasks/:id", middlewares.AuthMiddleware(),controllers.DeleteTask)
	router.POST("/register", middlewares.AuthMiddleware(),controllers.CreateUser)
	router.GET("/users/:id", middlewares.AuthMiddleware(),controllers.GetUserById)
	router.POST("/login", middlewares.AuthMiddleware(),controllers.LogIn)
	router.GET("/users", middlewares.AuthMiddlewareRole([]string{"Admin", "Super Admin"}), controllers.GetAllUsers)
	router.GET("/me", middlewares.AuthMiddleware(), controllers.Me)
	router.PUT("/users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
	router.DELETE("/users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)


	router.Run()
}