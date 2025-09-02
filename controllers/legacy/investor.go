package legacy

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"proto/backendAPI/initializers"
	"proto/backendAPI/models"
)

func CreateInvestorFromLegacy(c *gin.Context) {
	var investorLegacy models.InvestorLegacy

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
