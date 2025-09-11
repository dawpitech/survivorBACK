package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"gorm.io/gorm"
	"reflect"
)

func GetAllFounders(_ *gin.Context) (*[]models.Founder, error) {
	var founders []models.Founder
	if result := initializers.DB.Find(&founders); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &founders, nil
}

func GetFounder(_ *gin.Context, in *routes.GetFounderRequest) (*models.Founder, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var founder models.Founder
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("Startup").Find(&founder); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Founder not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	return &founder, nil
}

func CreateNewFounder(c *gin.Context, in *routes.FounderCreationRequest) (*models.Founder, error) {
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

	founder := models.Founder{
		UUID:        uuid.New().String(),
		ID:          nil,
		Name:        in.Name,
		StartupUUID: nil,
		StartupID:   nil,
		Startup:     nil,
	}

	if err := initializers.DB.Create(&founder); err.Error != nil {
		return nil, errors.New("Internal server error")
	}

	return &founder, nil
}

func DeleteFounder(c *gin.Context, in *routes.DeleteStartupRequest) error {
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

	var founder models.Founder
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&founder); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "Founder not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if rst := initializers.DB.Delete(&founder); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func UpdateFounder(c *gin.Context, in *routes.FounderUpdateRequest) (*models.Founder, error) {
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

	if authUser.Role != "admin" && authUser.UUID != in.UUID {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var founder models.Founder
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("Startups").First(&founder); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Founder not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

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

	if err := initializers.DB.Model(&founder).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &founder, nil
}
