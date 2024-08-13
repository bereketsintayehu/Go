package controllers

import (
	"net/http"
	"task-manager/domain"
	"task-manager/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userUseCase *usecase.UserUseCase
}

func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	userRole := c.MustGet("role").(string)

	users, err := uc.userUseCase.GetAllUsers(userRole)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userRole := c.MustGet("role").(string)
	userId := c.MustGet("userId").(primitive.ObjectID)

	user, err := uc.userUseCase.GetUserByID(userRole, userId, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) Me(c *gin.Context) {
	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)

	user, err := uc.userUseCase.GetUserByID(userRole, userId, userId.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
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

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)
	userRole := c.MustGet("role").(string)

	updatedUser, err := uc.userUseCase.UpdateUser(userRole, userId.Hex(), objectId.Hex(), &user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": updatedUser})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	deletedUser, err := uc.userUseCase.DeleteUser(c.MustGet("role").(string), c.MustGet("userId").(primitive.ObjectID).Hex(), id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "user": deletedUser})
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Role.String() == "" {
		user.Role = domain.UserRole(0)
	}

	// role if its in context, else default to User
	userRole, roleExists := c.Get("role")
	if !roleExists {
		userRole = "User"
	} 

	createdUser, err := uc.userUseCase.CreateUser(userRole.(string), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": createdUser})
}

func (uc *UserController) LogIn(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := uc.userUseCase.Login(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
