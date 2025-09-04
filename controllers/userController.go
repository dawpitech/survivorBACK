package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func GetMe(c *gin.Context) {
	userInterface, exist := c.Get("currentUser")

	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var user models.User
	switch u := userInterface.(type) {
	case models.User:
		user = u
	case *models.User:
		user = *u
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user.GetPublicUser())
}

func GetUser(c *gin.Context) {
	uuidParam := c.Param("uuid")

	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", uuidParam).Find(&userFound); rst.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown UUID"})
		return
	}

	c.JSON(http.StatusOK, userFound.GetPublicUser())
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

func DeleteUser(c *gin.Context) {
	uuidParam := c.Param("uuid")

	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", uuidParam).Find(&userFound); rst.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown UUID"})
		return
	}

	if rst := initializers.DB.Delete(&userFound); rst.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func UpdateUser(c *gin.Context) {
	uuidParam := c.Param("uuid")

	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var userFound models.User
	if rst := initializers.DB.Where("uuid=?", uuidParam).Find(&userFound); rst.Error != nil {
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
	delete(updates, "founder_uuid")
	delete(updates, "founder_id")
	delete(updates, "investor_uuid")
	delete(updates, "investor_id")
	delete(updates, "role")

	if val, ok := updates["password"]; ok {
		if passwordStr, ok := val.(string); ok && passwordStr != "" {
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			updates["password"] = string(passwordHash)
		}
	}

	if err := initializers.DB.Model(&userFound).Updates(updates); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
