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

func GetAllNews(_ *gin.Context) (*[]models.NewsDetails, error) {
	var news []models.NewsDetails
	if result := initializers.DB.Find(&news); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &news, nil
}

func GetNews(_ *gin.Context, in *routes.GetNewsRequest) (*models.NewsDetails, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var news models.NewsDetails
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "News not found")
		} else {
			return nil, errors.New("Internal server error")
		}
	}

	return &news, nil
}

func CreateNewNews(_ *gin.Context, in *routes.NewsCreationRequest) (*models.NewsDetails, error) {
	news := models.NewsDetails{
		News: models.News{
			UUID:      uuid.New().String(),
			ID:        nil,
			Location:  nil,
			Title:     in.Title,
			Category:  nil,
			StartupId: nil,
		},
		Description: "",
		NewsPicture: nil,
	}

	if createResult := initializers.DB.Create(&news); createResult.Error != nil {
		return nil, errors.New("Internal server error")
	}

	return &news, nil
}

func DeleteNews(_ *gin.Context, in *routes.DeleteNewsRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var news models.News
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewUserNotFound(nil, "News not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if rst := initializers.DB.Delete(&news); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func UpdateNews(_ *gin.Context, in *routes.NewsUpdateRequest) (*models.NewsDetails, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var news models.NewsDetails
	if rst := initializers.DB.Where("uuid=?", in.UUID).First(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewUserNotFound(nil, "News not found")
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

	if err := initializers.DB.Model(&news).Updates(updates).Error; err != nil {
		return nil, errors.New("Internal server error")
	}

	return &news, nil
}

func GetNewsPicture(c *gin.Context, in *routes.GetNewsPictureRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var news models.NewsDetails
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("NewsPicture").First(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "News not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	if news.NewsPicture == nil || len(news.NewsPicture.Picture) == 0 {
		return errors.NewNotFound(nil, "News picture not found")
	}

	picture := news.NewsPicture.Picture

	c.Data(http.StatusOK, "image/png", picture)
	return nil
}

func UpdateNewsPicture(c *gin.Context) error {
	userUUID := c.Param("uuid")
	file, err := c.FormFile("picture")

	if userUUID == "" {
		return errors.NewNotFound(nil, "News not found")
	}

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Internal server error")
	}

	var news models.NewsDetails
	if rst := initializers.DB.Where("uuid=?", userUUID).Preload("NewsPicture").First(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "News not found")
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

	newsPicture := models.NewsPicture{
		NewsUUID: news.UUID,
		Picture:  fileBytes,
	}

	if rst := initializers.DB.Save(&newsPicture); rst.Error != nil {
		return errors.New("Internal server error")
	}
	return nil
}

func ResetNewsPicture(_ *gin.Context, in *routes.ResetNewsPictureRequest) error {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return errors.NewNotValid(nil, "Invalid UUID")
	}

	var news models.NewsDetails
	if rst := initializers.DB.Where("uuid=?", in.UUID).Preload("UserPicture").First(&news); rst.Error != nil {
		if errors.Is(rst.Error, gorm.ErrRecordNotFound) {
			return errors.NewNotFound(nil, "News not found")
		} else {
			return errors.New("Internal server error")
		}
	}

	utils.ResetNewsPicture(&news)
	return nil
}
