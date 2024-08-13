package repository

import (
	"task-manager/domain"
	"context"
	"log"
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &UserRepository{collection: collection}
}

func (r *UserRepository) EnsureSuperAdmin() error{
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    email := os.Getenv("SUPER_ADMIN_EMAIL")
    password := os.Getenv("SUPER_ADMIN_PASSWORD")
    if email == "" || password == "" {
        log.Fatal("SUPER_ADMIN_EMAIL or SUPER_ADMIN_PASSWORD environment variable not set")
    }

    filter := bson.M{"email": email}
    user := domain.User{}
    err = r.collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        user = domain.User{
            Email:    email,
            Password: password,
            Role:     domain.Super,
        }
        _, err = r.collection.InsertOne(context.TODO(), user)
        if err != nil {
            log.Fatal(err)
        }
        log.Println("Super admin created")
    }
	return nil
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	users := []domain.User{}
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		user := domain.User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) FindByID(id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}
	user := domain.User{}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	filter := bson.M{"email": email}
	user := domain.User{}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAllAdmin() ([]domain.User, error) {
	users := []domain.User{}
	cursor, err := r.collection.Find(context.TODO(), bson.M{
		"role": bson.M{
			"$in": []string{fmt.Sprint(domain.Admin), fmt.Sprint(domain.UserR)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		user := domain.User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepository) CreateUser(user domain.User) (*domain.User, error) {
	mu.Lock()
	defer mu.Unlock()

	res, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *UserRepository) DeleteUser(id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}
	user := domain.User{}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user domain.User) (*domain.User, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"email": user.Email,
			"password": user.Password,
			"role": user.Role,
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	updatedUser := domain.User{}
	err = r.collection.FindOne(context.TODO(), filter).Decode(&updatedUser)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &updatedUser, nil
}
