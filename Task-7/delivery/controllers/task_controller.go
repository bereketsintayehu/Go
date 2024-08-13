package controllers

import (
	"net/http"
	"task-manager/domain"
	"task-manager/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUseCase *usecase.TaskUseCase
}

func NewTaskController(taskUseCase *usecase.TaskUseCase) *TaskController {
	return &TaskController{TaskUseCase: taskUseCase}
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userId").(primitive.ObjectID)

	tasks, err := tc.TaskUseCase.GetAllTasks(userRole, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	userRole := c.MustGet("role").(string)
	userId := c.MustGet("userId").(primitive.ObjectID)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}
	task, err := tc.TaskUseCase.GetTaskByID(userRole, userId, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userId").(primitive.ObjectID)

	task := domain.Task{}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTask, err := tc.TaskUseCase.CreateTask(userRole, userID, &task)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": createdTask})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userId").(primitive.ObjectID)

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask, err := tc.TaskUseCase.UpdateTask(userRole, userID, objectId, &task)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": updatedTask})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userId").(primitive.ObjectID)

	deletedTask, err := tc.TaskUseCase.DeleteTask(userRole, userID, objectId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": deletedTask})
}