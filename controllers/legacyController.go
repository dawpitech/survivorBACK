package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/legacy"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func CreateUserFromLegacy(c *gin.Context) {
	var userLegacy legacy.UserLegacy

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

	if createResult := initializers.DB.Create(&user); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func CreateInvestorFromLegacy(c *gin.Context) {
	var investorLegacy legacy.InvestorLegacy

	if err := c.ShouldBindJSON(&investorLegacy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if investorLegacy.Email == "" || investorLegacy.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more required field empty."})
		return
	}

	investor := models.Investor{
		UUID:            uuid.New().String(),
		ID:              &investorLegacy.ID,
		Name:            investorLegacy.Name,
		LegalStatus:     investorLegacy.LegalStatus,
		Address:         investorLegacy.Address,
		Email:           investorLegacy.Email,
		Phone:           investorLegacy.Phone,
		CreatedAt:       investorLegacy.CreatedAt,
		Description:     investorLegacy.Description,
		InvestorType:    investorLegacy.InvestorType,
		InvestmentFocus: investorLegacy.InvestmentFocus,
	}

	if createResult := initializers.DB.Create(&investor); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
