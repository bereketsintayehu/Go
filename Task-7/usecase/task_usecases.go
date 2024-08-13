package usecase

import (
	"errors"
	"task-manager/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) GetAllTasks(userRole string, userID primitive.ObjectID) ([]domain.Task, error) {
	if userRole == "User" {
		return uc.GetTasksByUser(userID)
	}
	return uc.repo.FindAll()
}

func (uc *TaskUseCase) GetTasksByUser(userId primitive.ObjectID) ([]domain.Task, error) {
	return uc.repo.FindByUserID(userId)
}

func (uc *TaskUseCase) GetTaskByID(userRole string, userId, id primitive.ObjectID) (*domain.Task, error) {

	task, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if userRole == "User" && task.CreatedBy != userId {
		return nil, errors.New("unauthorized")
	}
	return task, nil
}

func (uc *TaskUseCase) CreateTask(userRole string, userID primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
	if userRole == "User" {
		task.CreatedBy = userID
	}

	return uc.repo.Create(task)
}

func (uc *TaskUseCase) UpdateTask(userRole string, userID, id primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
	taskToBeUpdated, err := uc.GetTaskByID(userRole, userID, id)
	if err != nil {
		return nil, err
	}
	if userRole == "User" && taskToBeUpdated.CreatedBy != userID {
		return nil, errors.New("unauthorized")
	}

	_, err = uc.repo.Update(id, task)
	if err != nil {
		return nil, err
	}
	return uc.GetTaskByID(userRole, userID, id)
}

func (uc *TaskUseCase) DeleteTask(userRole string, userID, id primitive.ObjectID) (*domain.Task, error) {
	taskToBeDeleted, err := uc.GetTaskByID(userRole, userID, id)
	if err != nil {
		return nil, err
	}
	if userRole == "User" && taskToBeDeleted.CreatedBy != userID {
		return nil, errors.New("unauthorized")
	}
	return uc.repo.Delete(id)
}
