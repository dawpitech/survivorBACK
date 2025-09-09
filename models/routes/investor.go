package routes

type GetInvestorRequest = GenericUUIDFromPath
type DeleteInvestorRequest = GenericUUIDFromPath

type InvestorCreationRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type InvestorUpdateRequest struct {
	GenericUUIDFromPath
	Name            *string `json:"name"`
	LegalStatus     *string `json:"legal_status"`
	Address         *string `json:"address"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	Description     *string `json:"description"`
	InvestorType    *string `json:"investor_type"`
	InvestmentFocus *string `json:"investment_focus"`
}
