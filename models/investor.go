package models

type Investor struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"id" gorm:"unique;index"` // legacy

	Name            string  `json:"name"`
	LegalStatus     *string `json:"legal_status"`
	Address         *string `json:"address"`
	Email           string  `json:"email" gorm:"unique;not null;default:null"`
	Phone           *string `json:"phone"`
	CreatedAt       *string `json:"created_at"`
	Description     *string `json:"description"`
	InvestorType    *string `json:"investor_type"`
	InvestmentFocus *string `json:"investment_focus"`

	InvestorPicture *InvestorPicture `json:"investor_picture" gorm:"foreignKey:InvestorUUID;references:UUID"`
}

type InvestorPicture struct {
	InvestorUUID string `json:"investor_uuid" gorm:"type:uuid;primary_key"`
	Picture  []byte `json:"picture" gorm:"type:bytea"`
}
