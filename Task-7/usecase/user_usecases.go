package usecase

import (
	"errors"
	"task-manager/domain"
	"task-manager/infrastructure"
	"task-manager/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	userRepo        *repository.UserRepository
	passwordService *infrastructure.PasswordService
	jwtService      *infrastructure.JWTService
}

func NewUserUseCase(userRepo *repository.UserRepository, passwordService *infrastructure.PasswordService, jwtService *infrastructure.JWTService) domain.UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *UserUseCase) CreateUser(userRole string, user *domain.User) (*domain.User, error) {
	if userRole == "User" && user.Role.String() != "User" {
		return nil, errors.New("unauthorized")
	}

	if userRole == "Admin" && user.Role.String() == "SuperAdmin" {
		return nil, errors.New("unauthorized")
	}

	hashedPassword, err := uc.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()

	createdUser, err := uc.userRepo.CreateUser(*user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (uc *UserUseCase) GetUserByID(userRole string, userId primitive.ObjectID, id string) (*domain.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if userRole == "User" && id != userRole {
		return nil, errors.New("unauthorized")
	}

	user, err := uc.userRepo.FindByID(objectId)
	if err != nil {
		return nil, err
	}
	return user, nil

}


func (uc *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) AuthenticateUser(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	if !uc.passwordService.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	jwtToken, err := uc.jwtService.GenerateToken(user.ID.Hex(), user.Role.String())
	if err != nil {
		return "", errors.New("internal server error")
	}

	return jwtToken, nil
}

func (uc *UserUseCase) GetAllUsers(userRole string) ([]domain.User, error) {
	if userRole == "SuperAdmin" {
		return uc.userRepo.FindAll()
	} else if userRole == "Admin" {
		return uc.userRepo.FindAllAdmin()
	}
	return nil, errors.New("unauthorized")
}

func (uc *UserUseCase) DeleteUser(userRole, userId, userToBeDeletedID string) (*domain.User, error) {
	_, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	toBeDeletedID, err := primitive.ObjectIDFromHex(userToBeDeletedID)
	if err != nil {
		return nil, err
	}

	userToBeDeleted, err := uc.userRepo.FindByID(toBeDeletedID)
	if err != nil {
		return nil, err
	}

	if userToBeDeleted == nil {
		return nil, errors.New("user not found")
	}

	if userRole == "Admin"  && userToBeDeleted.Role.String() == "SuperAdmin" {
		return nil, errors.New("unauthorized")
	}

	if userRole == "User" {
		if userId != userToBeDeletedID {
			return nil, errors.New("unauthorized")
		}
	}

	return uc.userRepo.DeleteUser(toBeDeletedID)

}

func (uc *UserUseCase) UpdateUser(userRole, userId, userToBeUpdatedID string, user *domain.User) (*domain.User, error) {
	_, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	toBeUpdatedID, err := primitive.ObjectIDFromHex(userToBeUpdatedID)
	if err != nil {
		return nil, err
	}

	userToBeUpdated, err := uc.userRepo.FindByID(toBeUpdatedID)
	if err != nil {
		return nil, err
	}

	if userRole == "User" && userId != userToBeUpdatedID {
		return nil, errors.New("unauthorized")
	}

	if userRole == "Admin" && userToBeUpdated.Role.String() == "SuperAdmin" {
		return nil, errors.New("unauthorized")
	}

	updatedUser, err := uc.userRepo.UpdateUser(*user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
func (uc *UserUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !uc.passwordService.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := uc.jwtService.GenerateToken(user.ID.Hex(), user.Role.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
