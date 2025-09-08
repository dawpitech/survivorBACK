package models

type Partner struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"-" gorm:"unique;index"` // legacy

	Name            string  `json:"name"`
	LegalStatus     *string `json:"legal_status"`
	Address         *string `json:"address"`
	Email           string  `json:"email" gorm:"unique;not null;default:null"`
	Phone           *string `json:"phone"`
	CreatedAt       *string `json:"created_at"`
	Description     *string `json:"description"`
	PartnershipType *string `json:"partnership_type"`
}
