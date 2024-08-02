package data

import (
	"log"
	"sync"
	"fmt"
	"task-manager/db"
	"task-manager/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"

)

var (
	client *mongo.Client
	collection *mongo.Collection
	mu sync.Mutex
)

func init() {
	var err error
	client, err = db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("taskdb").Collection("task-manager")
}

func GetAllTasks() []models.Task {
	tasks := []models.Task{}
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		task := models.Task{}
		err := cursor.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func GetTaskById(id primitive.ObjectID) *models.Task {
	filter := bson.M{"_id": id}
	task := models.Task{}
	err := collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return nil
	}
	return &task
}

func CreateTask(task models.Task) *models.Task {
	mu.Lock()
	defer mu.Unlock()
	
	insertResult, err := collection.InsertOne(context.TODO(), task)
    if err != nil {
        log.Fatal(err)
    }

    task.ID = insertResult.InsertedID.(primitive.ObjectID)

    fmt.Println(insertResult)

    return &task
}

func UpdateTask(id primitive.ObjectID, task models.Task) *mongo.UpdateResult {
	mu.Lock()
	defer mu.Unlock()
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{
		"title": task.Title,
		"description": task.Description,
		"status": task.Status,
	}}
	UpdateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return UpdateResult
}

func DeleteTask(id primitive.ObjectID) *models.Task {
	mu.Lock()
	defer mu.Unlock()
	filter := bson.M{"id": id}
	task := GetTaskById(id)
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return task
}
