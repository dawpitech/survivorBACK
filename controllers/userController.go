package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(_ *gin.Context, _ *struct{}) (*[]models.PublicUser, error) {
	var users []models.User
	var publicUsers []models.PublicUser
	if result := initializers.DB.Find(&users); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	for _, user := range users {
		publicUsers = append(publicUsers, user.GetPublicUser())
	}
	return &publicUsers, nil
}

func GetMe(c *gin.Context, _ *struct{}) (*models.PublicUser, error) {
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var user models.User
	switch u := userInterface.(type) {
	case models.User:
		user = u
	case *models.User:
		user = *u
	default:
		return nil, errors.New("Internal server error")
	}

	userPublic := user.GetPublicUser()
	return &userPublic, nil
}

func GetUser(_ *gin.Context, in *routes.GetUserRequest) (*models.PublicUser, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var user models.User
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&user); rst.Error != nil {
		return nil, errors.NewUserNotFound(nil, "User not found")
	}

	userFoundPublic := user.GetPublicUser()
	return &userFoundPublic, nil
}

func CreateNewUser(_ *gin.Context, in *routes.UserCreationRequest) (*models.PublicUser, error) {
	var userFound models.User
	if findResult := initializers.DB.Where("email=?", in.Email).Find(&userFound); findResult.Error != nil {
		return nil, errors.NewUserNotFound(nil, "User not found")
	}

	if userFound.Password != nil && userFound.UUID != "" {
		return nil, errors.NewAlreadyExists(nil, "Email already used")
	}

	user := models.User{
		UUID:         uuid.New().String(),
		ID:           nil,
		Name:         in.Name,
		Email:        in.Email,
		Password:     nil,
		Role:         in.Role,
		FounderUUID:  nil,
		FounderID:    nil,
		InvestorUUID: nil,
		InvestorID:   nil,
	}

	if createResult := initializers.DB.Create(&user); createResult.Error != nil {
		return nil, errors.New("Internal server error")
	}

	publicUser := user.GetPublicUser()
	return &publicUser, nil
}

func DeleteUser(_ *gin.Context, in *routes.DeleteUserRequest) (*struct{}, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&userFound); rst.Error != nil {
		return nil, errors.NewUserNotFound(nil, "User not found")
	}

	if rst := initializers.DB.Delete(&userFound); rst.Error != nil {
		return nil, errors.New("Internal server error")
	}

	var empty struct{}
	return &empty, nil
}

func UpdateUser(_ *gin.Context, in *routes.UpdateUserRequest) (*models.PublicUser, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid = ?", in.UUID).First(&userFound); rst.Error != nil {
		return nil, errors.NewUserNotFound(nil, "User not found")
	}

	if in.Name == "" && in.Email == "" && in.Password == "" {
		return nil, errors.NewNotValid(nil, "Invalid elements")
	}

	updates := make(map[string]interface{})
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Email != "" {
		updates["email"] = in.Email
	}
	if in.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(passwordHash)
	}

	if err := initializers.DB.Model(&userFound).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	publicUser := userFound.GetPublicUser()
	return &publicUser, nil
}
