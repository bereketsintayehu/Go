package domain

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
	Cancelled
)

func (ts TaskStatus) String() string {
	return [...]string{"Pending", "InProgress", "Completed", "Cancelled"}[ts]
}

func (ts TaskStatus) MarshalText() ([]byte, error) {
	return []byte(ts.String()), nil
}

func (ts TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}

func (ts *TaskStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	switch statusStr {
	case "Pending":
		*ts = Pending
	case "InProgress":
		*ts = InProgress
	case "Completed":
		*ts = Completed
	case "Cancelled":
		*ts = Cancelled
	default:
		return errors.New("invalid TaskStatus")
	}

	return nil
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      TaskStatus         `json:"status"`
	CreatedBy   primitive.ObjectID `json:"created_by"`
}

type TaskRepository interface {
	FindAll() ([]Task, error)
	FindByID(id primitive.ObjectID) (*Task, error)
	FindByUserID(userID primitive.ObjectID) ([]Task, error)
	Create(task *Task) (*Task, error)
	Update(id primitive.ObjectID,task *Task) (*mongo.UpdateResult, error)
	Delete(id primitive.ObjectID) (*Task, error)
}

type TaskUseCase interface {
	GetAllTasks(userRole string, userID primitive.ObjectID) ([]Task, error) 
	GetTasksByUser(userId primitive.ObjectID) ([]Task, error)
	GetTaskByID(userRole string, userId, id primitive.ObjectID) (*Task, error)
	CreateTask(userRole string, userID primitive.ObjectID, task *Task) (*Task, error)
	UpdateTask(userRole string, userID, id primitive.ObjectID, task *Task) (*Task, error)
	DeleteTask(userRole string, userID, id primitive.ObjectID) (*Task, error)
}
