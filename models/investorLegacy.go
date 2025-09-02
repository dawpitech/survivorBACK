package models

type InvestorLegacy struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	LegalStatus     *string `json:"legal_status"`
	Address         *string `json:"address"`
	Email           string  `json:"email"`
	Phone           *string `json:"phone"`
	CreatedAt       *string `json:"created_at"`
	Description     *string `json:"description"`
	InvestorType    *string `json:"investor_type"`
	InvestmentFocus *string `json:"investment_focus"`
}
