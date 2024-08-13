package repository

import (
	"task-manager/domain"
	"context"
	"log"
	"sync"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	mu sync.Mutex
)

type TaskRepository struct {
    collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) domain.TaskRepository {
	return &TaskRepository{collection: collection}
}

func (r *TaskRepository) FindAll() ([]domain.Task, error) {
	tasks := []domain.Task{}
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		task := domain.Task{}
		err := cursor.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) FindByUserID(userId primitive.ObjectID) ([]domain.Task, error) {
	tasks := []domain.Task{}
	filter := bson.M{"createdBy": userId}
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		task := domain.Task{}
		err := cursor.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) FindByID(id primitive.ObjectID) (*domain.Task, error) {
	filter := bson.M{"_id": id}
	task := domain.Task{}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Create(task *domain.Task) (*domain.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	
	insertResult, err := r.collection.InsertOne(context.TODO(), task)
    if err != nil {
        log.Fatal(err)
    }

    task.ID = insertResult.InsertedID.(primitive.ObjectID)

    return task, nil
}

func (r *TaskRepository) Update(id primitive.ObjectID, task *domain.Task) (*mongo.UpdateResult, error) {
	mu.Lock()
	defer mu.Unlock()
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{
		"title": task.Title,
		"description": task.Description,
		"status": task.Status,
	}}
	UpdateResult, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return UpdateResult, nil
}

func (r *TaskRepository) Delete(id primitive.ObjectID) (*domain.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	filter := bson.M{"id": id}
	task := domain.Task{}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return &task, nil
}
