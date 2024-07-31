package data

import (
	"sync"
	"task-manager/models"
)

var (
	tasks = []models.Task{
		{ID: "1", Title: "Task 1", Description: "This is task 1", Status: models.InProgress},
		{ID: "2", Title: "Task 2", Description: "This is task 2", Status: models.Completed},
		{ID: "3", Title: "Task 3", Description: "This is task 3", Status: models.Pending},
	}
	mu sync.Mutex
)

func GetAllTasks() []models.Task {
	return tasks
}

func GetTaskById(id string) *models.Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func CreateTask(task models.Task) *models.Task {
	mu.Lock()
	defer mu.Unlock()

	tasks = append(tasks, task)
	return &task
}

func UpdateTask(id string, task models.Task) {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			tasks[i] = task
		}
	}
}

func DeleteTask(id string) *models.Task {
	mu.Lock()
	defer mu.Unlock()
	deletedTask := models.Task{}
	for i, t := range tasks {
		if t.ID == id {
			deletedTask = tasks[i]
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}
	return &deletedTask
}
