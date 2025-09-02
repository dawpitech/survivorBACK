package models

type UserLegacy struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	FounderID  *uint  `json:"founder_id"`
	InvestorID *uint  `json:"investor_id"`
	ID         *uint  `json:"id"`
}
