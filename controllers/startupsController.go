package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

func GetAllStartups(_ *gin.Context, _ *struct{}) (*[]models.StartupDetail, error) {
	var startups []models.StartupDetail
	if result := initializers.DB.Find(&startups); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &startups, nil
}

func GetStartup(_ *gin.Context, in *routes.GetStartupRequest) (*models.StartupDetail, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var startup models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("Founders").Find(&startup); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Startup not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	return &startup, nil
}

func CreateNewStartup(c *gin.Context, in *routes.StartupCreationRequest) (*models.StartupDetail, error) {
	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	if authUser.Role != "admin" {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("email=?", in.Email).Find(&startupFound); rst.Error == nil {
		return nil, errors.NewAlreadyExists(nil, "Email already used")
	} else {
		if !errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Internal server error")
		}
	}

	currentDate := time.Now().Format("2006-01-02")
	startup := models.StartupDetail{
		StartupList: models.StartupList{
			UUID:        uuid.New().String(),
			ID:          nil,
			Name:        in.Name,
			LegalStatus: nil,
			Address:     nil,
			Email:       in.Email,
			Phone:       nil,
			Sector:      nil,
			Maturity:    nil,
		},
		CreatedAt:      &currentDate,
		Description:    nil,
		WebsiteUrl:     nil,
		SocialMediaURL: nil,
		ProjectStatus:  nil,
		Needs:          nil,
		Founders:       nil,
	}

	if err := initializers.DB.Create(&startup); err.Error != nil {
		return nil, errors.New("Internal server error")
	}

	return &startup, nil
}

func DeleteStartup(c *gin.Context, in *routes.DeleteStartupRequest) error {
	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return errors.New("Internal server error")
	}

	if authUser.Role != "admin" {
		return errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&startupFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "Startup not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if rst := initializers.DB.Delete(&startupFound); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func UpdateStartup(c *gin.Context, in *routes.UpdateStartupRequest) (*models.StartupDetail, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("Founders").First(&startupFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Startup not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	if startupFound.UUID == "" {
		return nil, errors.NewUserNotFound(nil, "Startup not found")
	}

	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	isFounder := false
	for _, f := range startupFound.Founders {
		if authUser.UUID == f.UUID {
			isFounder = true
			break
		}
	}
	if authUser.Role != "admin" && !isFounder {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	updates := make(map[string]interface{})

	val := reflect.ValueOf(*in)
	typ := reflect.TypeOf(*in)
	hasUpdate := false

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" || jsonTag == "uuid" {
			continue
		}
		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			hasUpdate = true
			updates[jsonTag] = fieldValue.String()
		}
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			strVal, ok := fieldValue.Elem().Interface().(string)
			if ok && strVal != "" {
				hasUpdate = true
				updates[jsonTag] = strVal
			}
		}
	}

	if !hasUpdate {
		return nil, errors.NewNotValid(nil, "Invalid body")
	}

	if err := initializers.DB.Model(&startupFound).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &startupFound, nil
}

func AddViewToStartup(_ *gin.Context, in *routes.UpdateStartupRequest) (*models.StartupDetail, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("Founders").First(&startupFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Startup not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	if startupFound.UUID == "" {
		return nil, errors.NewUserNotFound(nil, "Startup not found")
	}

	startupFound.ViewsCount += 1

	if err := initializers.DB.Save(&startupFound).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &startupFound, nil
}

func UploadStartupFile(c *gin.Context) error {
	startupUUID := c.Param("uuid")
	file, err := c.FormFile("file")
	if err != nil {
		return errors.NewBadRequest(err, "No file given")
	}

	if _, err := uuid.Parse(startupUUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}
	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", startupUUID).Preload("Founders").First(&startupFound); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "Startup not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if startupFound.UUID == "" {
		return errors.NewUserNotFound(nil, "Startup not found")
	}

	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return errors.New("Internal server error")
	}

	isFounder := false
	for _, f := range startupFound.Founders {
		if authUser.UUID == f.UUID {
			isFounder = true
			break
		}
	}
	if authUser.Role != "admin" && !isFounder {
		return errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	uploadDir := "./startup_files"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return errors.New("Internal server error")
	}
	filePath := filepath.Join(uploadDir, startupUUID)

	_ = os.Remove(filePath)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return errors.New("Internal server error")
	}

	return nil
}

func GetStartupFile(c *gin.Context) error {
	startupUUID := c.Param("uuid")
	filePath := filepath.Join("./startup_files", startupUUID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.NewNotFound(nil, "File not found")
	}
	c.File(filePath)
	return nil
}
