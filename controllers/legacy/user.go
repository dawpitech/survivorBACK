package legacy

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"proto/backendAPI/initializers"
	"proto/backendAPI/models"
)

func CreateUserFromLegacy(c *gin.Context) {
	var userLegacy models.UserLegacy

	if err := c.ShouldBindJSON(&userLegacy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userLegacy.Email == "" || userLegacy.Name == "" || userLegacy.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more required field empty."})
		return
	}

	user := models.User{
		UUID:         uuid.New().String(),
		ID:           &userLegacy.ID,
		Name:         userLegacy.Name,
		Email:        userLegacy.Email,
		Password:     nil,
		Role:         userLegacy.Role,
		FounderUUID:  nil,
		FounderID:    userLegacy.FounderID,
		InvestorUUID: nil,
		InvestorID:   userLegacy.InvestorID,
	}

	if user.FounderID != nil {
		founderUUID := uuid.New().String()
		user.FounderUUID = &founderUUID
	}

	if user.InvestorID != nil {
		investorUUID := uuid.New().String()
		user.InvestorUUID = &investorUUID
	}

	if createResult := initializers.DB.Create(&user); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
