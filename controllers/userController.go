package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if result := initializers.DB.Find(&users); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateNewUser(c *gin.Context) {
	var userCreationRequest models.UserCreationRequest

	if err := c.ShouldBindJSON(&userCreationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	if findResult := initializers.DB.Where("email=?", userCreationRequest.Email).Find(&userFound); findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if userFound.Password != nil && userFound.UUID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already used"})
		return
	}

	user := models.User{
		UUID:         uuid.New().String(),
		ID:           nil,
		Name:         userCreationRequest.Name,
		Email:        userCreationRequest.Email,
		Password:     nil,
		Role:         userCreationRequest.Role,
		FounderUUID:  nil,
		FounderID:    nil,
		InvestorUUID: nil,
		InvestorID:   nil,
	}

	if createResult := initializers.DB.Create(&user); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
