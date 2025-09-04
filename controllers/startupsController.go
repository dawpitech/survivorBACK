package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GetAllStartups(c *gin.Context) {
	var startups []models.StartupDetail
	if result := initializers.DB.Find(&startups); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, startups)
}

func GetStartup(c *gin.Context) {
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

	c.JSON(http.StatusOK, startupFound)
}

func CreateNewStartup(c *gin.Context) {
	var startupCreationRequest models.StartupCreationRequest

	if err := c.ShouldBindJSON(&startupCreationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startupFound models.StartupDetail
	if findResult := initializers.DB.Where("email=?", startupFound.Email).Find(&startupFound); findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if startupFound.UUID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already used"})
		return
	}

	currentDate := time.Now().Format("2006-01-02")
	startup := models.StartupDetail{
		StartupList: models.StartupList{
			UUID:        uuid.New().String(),
			ID:          nil,
			Name:        startupCreationRequest.Name,
			LegalStatus: nil,
			Address:     nil,
			Email:       startupCreationRequest.Email,
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

	if createResult := initializers.DB.Create(&startup); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
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
