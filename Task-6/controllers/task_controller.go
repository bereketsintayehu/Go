package controllers

import (
	"net/http"
	"task-manager/data"
	"task-manager/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTaks(c *gin.Context){
	userRole := c.MustGet("role").(string)

	if userRole == "User" {
		userId := c.MustGet("userId").(primitive.ObjectID)
		c.JSON(http.StatusOK, gin.H{"task": data.GetTasksByUserId(userId)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": data.GetAllTasks()})
}

func GetTaskById(c *gin.Context){
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
	
	task := data.GetTaskById(objectId)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func CreateTask(c *gin.Context){
	var task models.Task

	task.Status = models.Pending
	
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)
	task.CreatedBy = userId

	newTask := data.CreateTask(task)
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": newTask})
}

func UpdateTask(c *gin.Context){
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

	unapdatedTask := data.GetTaskById(objectId)
	if unapdatedTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)

	if userRole == "User" && unapdatedTask.CreatedBy != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not found."})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}



	data.UpdateTask(objectId, task)

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

func DeleteTask(c *gin.Context){
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

	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)
	taskToBeDeleted := data.GetTaskById(objectId)

	if userRole == "User" && taskToBeDeleted.CreatedBy != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not found."})
		return
	}
	deletedTask := data.DeleteTask(objectId)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": deletedTask})
}
