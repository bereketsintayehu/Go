package controllers

import (
	"net/http"
	"os"
	"strconv"
	"task-manager/data"
	"task-manager/models"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret []byte
	expirationTimeInt int
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if jwtSecret == nil {
		panic("JWT_SECRET environment variable not set")
	}

		expirationTimeStr := os.Getenv("JWT_EXPIRATION_TIME")
		expirationTimeInt, err = strconv.Atoi(expirationTimeStr)
		if err != nil {
			expirationTimeInt = 24
		}
}

func GetAllUsers(c *gin.Context) {
	userRole := c.MustGet("role").(string)
	if userRole == "Super Admin" {
		c.JSON(http.StatusOK, gin.H{"users": data.GetAllUsers()})
		return
	} else if userRole == "Admin" {
		c.JSON(http.StatusOK, gin.H{"users": data.GetAllUsersAdmin()})
		return
	}
	c.JSON(401, gin.H{"error": "Unauthorized"})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	user := data.GetUserById(objectId)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Me(c *gin.Context) {
	userId := c.MustGet("userId").(primitive.ObjectID)
	user := data.GetUserById(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	unapdatedUser := data.GetUserById(objectId)
	if unapdatedUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)

	if userRole == "User" && unapdatedUser.ID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found."})
		return
	}

	if userRole == "Admin" && unapdatedUser.Role == models.Super {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found."})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Role != 0 && userRole != "Super Admin" && user.Role != unapdatedUser.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can't update the role"})
	}

	data.UpdateUser(user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)
	userToBeDeleted := data.GetUserById(objectId)

	if userRole == "User" && userToBeDeleted.ID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found."})
		return
	}

	if userRole == "Admin" && userToBeDeleted.Role == models.Super {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found."})
		return
	}

	data.DeleteUser(objectId)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Role == 0 {
		user.Role = models.UserRole(0)
	}

	if user.Role != models.UserR {
		role := c.MustGet("role").(string)
		if role != "Super Admin" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user.Password = string(hashedPassword)
	newUser := data.CreateUser(user)
	if newUser == nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully", "user": newUser})
}

func LogIn(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	userByEmail := data.GetUserByEmail(user.Email)
	if userByEmail == nil {
		c.JSON(404, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userByEmail.Password), []byte(user.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	expirationTime := time.Now().Add(time.Duration(expirationTimeInt) * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userByEmail.ID,
		"email": userByEmail.Email,
		"role":  userByEmail.Role,
		"exp":   expirationTime,
	})

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User logged in successfully", "token": jwtToken})
}
