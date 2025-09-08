package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"FranceDeveloppe/JEB-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"net/http"
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
		return nil, errors.New("Internal server error")
	}

	if user.UUID == "" {
		return nil, errors.NewNotFound(nil, "User not found")
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

func DeleteUser(_ *gin.Context, in *routes.DeleteUserRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&userFound); rst.Error != nil {
		return errors.New("Internal server error")
	}

	if userFound.UUID == "" {
		return errors.NewUserNotFound(nil, "User not found")
	}

	if rst := initializers.DB.Delete(&userFound); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
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
		return nil, errors.NewNotValid(nil, "Invalid body")
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

func GetUserPicture(c *gin.Context, in *routes.GetUserPictureRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("UserPicture").First(&userFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "User not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if userFound.UserPicture == nil || len(userFound.UserPicture.Picture) == 0 {
		return errors.NewNotFound(nil, "User picture not found")
	}

	picture := userFound.UserPicture.Picture

	c.Data(http.StatusOK, "image/png", picture)
	return nil
}

func UpdateUserPicture(c *gin.Context) error {
	userUUID := c.Param("uuid")
	file, err := c.FormFile("picture")

	if userUUID == "" {
		return errors.NewNotFound(nil, "User not found")
	}

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Internal server error")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", userUUID).Preload("UserPicture").First(&userFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "User not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	openFile, openErr := file.Open()
	if openErr != nil {
		return errors.New("Internal server error")
	}
	defer func() { _ = openFile.Close() }()

	fileBytes, readErr := io.ReadAll(openFile)
	if readErr != nil {
		return errors.New("Internal server error")
	}

	userPicture := models.UserPicture{
		UserUUID: userFound.UUID,
		Picture:  fileBytes,
	}

	if rst := initializers.DB.Save(&userPicture); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func ResetUserPicture(_ *gin.Context, in *routes.ResetUserPictureRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("UserPicture").First(&userFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "User not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	utils.ResetUserPicture(&userFound)
	return nil
}
