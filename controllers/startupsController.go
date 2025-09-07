package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"net/http"
	"time"
)

func GetAllStartups(_ *gin.Context, _ *struct{}) (*[]models.StartupDetail, error) {
	var users []models.StartupDetail
	if result := initializers.DB.Find(&users); result.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &users, nil
}

func GetStartup(_ *gin.Context, in *routes.GetStartupRequest) (*models.StartupDetail, error) {
	if _, err := uuid.Parse(in.UUID); err != nil {
		return nil, errors.NewNotValid(nil, "Invalid UUID")
	}

	var startup models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", in.UUID).Find(&startup); rst.Error != nil {
		return nil, errors.NewUserNotFound(nil, "Startup not found")
	}

	return &startup, nil
}

func CreateNewStartup(_ *gin.Context, in *routes.StartupCreationRequest) (*models.StartupDetail, error) {
	var startupFound models.StartupDetail
	if findResult := initializers.DB.Where("email=?", in.Email).Find(&startupFound); findResult.Error != nil {
		return nil, errors.New("Internal server error")
	}

	if startupFound.UUID != "" {
		return nil, errors.NewAlreadyExists(nil, "Email already used")
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

func DeleteStartup(c *gin.Context) {
	uuidParam := c.Param("uuid")

	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", uuidParam).Find(&startupFound); rst.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown UUID"})
		return
	}

	if rst := initializers.DB.Delete(&startupFound); rst.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func UpdateStartup(c *gin.Context) {
	uuidParam := c.Param("uuid")

	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var startupFound models.StartupDetail
	if rst := initializers.DB.Where("uuid=?", uuidParam).Find(&startupFound); rst.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown UUID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	delete(updates, "uuid")
	delete(updates, "id")
	delete(updates, "created_at")

	if err := initializers.DB.Model(&startupFound).Updates(updates); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
