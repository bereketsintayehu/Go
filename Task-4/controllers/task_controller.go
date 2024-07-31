package controllers

import (
	"net/http"
	"task-manager/data"
	"task-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllTaks(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"task": data.GetAllTasks()})
}

func GetTaskById(c *gin.Context){
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	
	task := data.GetTaskById(id)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func CreateTask(c *gin.Context){
	var task models.Task

	// have default status to Pending
	task.Status = models.Pending
	
	task.ID = uuid.New().String()
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask := data.CreateTask(task)
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": newTask})
}

func UpdateTask(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}
	
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.ID != id  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update request"})
		return
	}
	

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

func DeleteTask(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	deletedTask := data.DeleteTask(id)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": deletedTask})
}
