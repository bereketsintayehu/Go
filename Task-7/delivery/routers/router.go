package routers

import (
	"task-manager/delivery/controllers"
	"task-manager/infrastructure"
	"task-manager/repository"
	"task-manager/usecase"
	"task-manager/bootstrap"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router    *gin.Engine
}

var (
	AuthMiddleware *infrastructure.AuthMiddleware
)
func (r *Router) userRoutes() {
	collection := bootstrap.GetCollecton("users")
	userRepository := repository.NewUserRepository(collection)
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService()
	userUseCase := usecase.NewUserUseCase(userRepository.(*repository.UserRepository), passwordService, jwtService)
	userController := controllers.NewUserController(userUseCase)
	r.Router.GET("/users", AuthMiddleware.AuthMiddlewareRole([]string{"Admin", "SuperAdmin"}), userController.GetAllUsers)
	r.Router.GET("/users/:id", AuthMiddleware.AuthMiddleware(), userController.GetUserById)
	r.Router.GET("/me", AuthMiddleware.AuthMiddleware(), userController.Me)
	r.Router.POST("/users", AuthMiddleware.AuthMiddleware(), userController.CreateUser)
	r.Router.POST("/login", userController.LogIn)
	r.Router.PUT("/users/:id", AuthMiddleware.AuthMiddleware(), userController.UpdateUser)
	r.Router.DELETE("/users/:id", AuthMiddleware.AuthMiddleware(), userController.DeleteUser)

}

func (r *Router) taskRoutes() {
	collection := bootstrap.GetCollecton("tasks")
	taskRepository := repository.NewTaskRepository(collection)
	taskUseCase := usecase.NewTaskUseCase(taskRepository.(*repository.TaskRepository))
	taskController := controllers.NewTaskController(taskUseCase)
	r.Router.GET("/tasks", AuthMiddleware.AuthMiddleware(), taskController.GetAllTasks)
	r.Router.GET("/tasks/:id", AuthMiddleware.AuthMiddleware(), taskController.GetTaskByID)
	r.Router.POST("/tasks", AuthMiddleware.AuthMiddleware(), taskController.CreateTask)
	r.Router.PUT("/tasks/:id", AuthMiddleware.AuthMiddleware(), taskController.UpdateTask)
	r.Router.DELETE("/tasks/:id", AuthMiddleware.AuthMiddleware(), taskController.DeleteTask)
}

func (r *Router) InitRoutes() {
	r.userRoutes()
	r.taskRoutes()
}

func NewRouter() *Router {
	r := gin.Default()
	return &Router{
		Router:    r,
	}

}