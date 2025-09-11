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
	"gorm.io/gorm"
	"io"
	"net/http"
	"reflect"
)

func GetAllEvents(_ *gin.Context) (*[]models.Event, error) {
	var events []models.Event
	if result := initializers.DB.Find(&events); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &events, nil
}

func GetEvent(_ *gin.Context, in *routes.GetEventRequest) (*models.Event, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Event not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	return &event, nil
}

func CreateNewEvent(c *gin.Context, in *routes.EventCreationRequest) (*models.Event, error) {
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

	event := models.Event{
		UUID:           uuid.New().String(),
		ID:             nil,
		Name:           in.Name,
		Date:           nil,
		Location:       nil,
		Description:    nil,
		EventType:      nil,
		TargetAudience: nil,
		EventPicture:   nil,
	}

	if createResult := initializers.DB.Create(&event); createResult.Error != nil {
		return nil, errors.New("Internal server error")
	}

	return &event, nil
}

func DeleteEvent(c *gin.Context, in *routes.DeleteUserRequest) error {
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

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "Event not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if rst := initializers.DB.Delete(&event); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func UpdateEvent(c *gin.Context, in *routes.UpdateEventRequest) (*models.Event, error) {
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

	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", in.UUID).First(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "Event not found")
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
			updates[jsonTag] = fieldValue.String()
			hasUpdate = true
		}
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			strVal, ok := fieldValue.Elem().Interface().(string)
			if ok && strVal != "" {
				updates[jsonTag] = strVal
				hasUpdate = true
			}
		}
	}

	if !hasUpdate {
		return nil, errors.NewNotValid(nil, "Invalid body")
	}

	if err := initializers.DB.Model(&event).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &event, nil
}

func GetEventPicture(c *gin.Context, in *routes.GetEventPictureRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("EventPicture").First(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "Event not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if event.EventPicture == nil || len(event.EventPicture.Picture) == 0 {
		return errors.NewNotFound(nil, "Event picture not found")
	}

	picture := event.EventPicture.Picture

	c.Data(http.StatusOK, "image/png", picture)
	return nil
}

func UpdateEventPicture(c *gin.Context) error {
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

	userUUID := c.Param("uuid")
	file, err := c.FormFile("picture")

	if userUUID == "" {
		return errors.NewNotFound(nil, "Event not found")
	}

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Internal server error")
	}

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", userUUID).Preload("EventPicture").First(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "Event not found")
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

	eventPicture := models.EventPicture{
		EventUUID: event.UUID,
		Picture:   fileBytes,
	}

	if rst := initializers.DB.Save(&eventPicture); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func ResetEventPicture(c *gin.Context, in *routes.ResetEventPictureRequest) error {
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

	var event models.Event
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("UserPicture").First(&event); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "Event not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	utils.ResetEventPicture(&event)
	return nil
}
