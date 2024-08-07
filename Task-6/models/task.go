package models

import (
	"encoding/json"
	"errors"
    "go.mongodb.org/mongo-driver/bson/primitive"
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

type Task struct{
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `json:"title"`
    Description string             `json:"description"`
    Status      TaskStatus         `json:"status"`
    CreatedBy   primitive.ObjectID `json:"created_by"`
}