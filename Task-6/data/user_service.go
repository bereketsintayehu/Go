package data

import (
	"context"
	"log"
	"os"
	"sync"
	"task-manager/db"
	"task-manager/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mu sync.Mutex
	client *mongo.Client
	collection *mongo.Collection
)

func init() {
	var err error
	client, err = db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("taskdb").Collection("users")

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("SUPER_ADMIN_EMAIL or SUPER_ADMIN_PASSWORD environment variable not set")
	}

	filter := bson.M{"email": email}
	user := models.User{}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		user = models.User{
			Email: email,
			Password: password,
			Role: models.Super,
		}
		log.Println("Super admin created")
		CreateUser(user)
	}

}

func GetAllUsers() []models.User {
	users := []models.User{}
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		user := models.User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}

func GetUserById(id primitive.ObjectID) *models.User {
	filter := bson.M{"_id": id}
	user := models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func GetUserByEmail(email string) *models.User {
	filter := bson.M{"email": email}
	user := models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func GetAllUsersAdmin() []models.User {
	users := []models.User{}
	cursor, err := collection.Find(context.TODO(), bson.M{
		"role": bson.M{
			"$in": []string{string(models.Admin), string(models.UserR)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		user := models.User{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}
func CreateUser(user models.User) *models.User {
	mu.Lock()
	defer mu.Unlock()

	res, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return &user
}

func DeleteUser(id primitive.ObjectID) {
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(user models.User) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"email": user.Email,
			"password": user.Password,
			"role": user.Role,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}
