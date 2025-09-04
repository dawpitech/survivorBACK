package legacy

type PartnerLegacy struct {
	ID              uint   `json:"id"`
	Name            string  `json:"name"`
	LegalStatus     *string `json:"legal_status"`
	Address         *string `json:"address"`
	Email           string  `json:"email"`
	Phone           *string `json:"phone"`
	CreatedAt       *string `json:"created_at"`
	Description     *string `json:"description"`
	PartnershipType *string `json:"partnership_type"`
}
