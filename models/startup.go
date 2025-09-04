package models

type StartupList struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"id" gorm:"unique;index"` // legacy

	Name        string  `json:"name"`
	LegalStatus *string `json:"legal_status"`
	Address     *string `json:"address"`
	Email       string  `json:"email" gorm:"unique;not null;default:null"`
	Phone       *string `json:"phone"`
	Sector      *string `json:"sector"`
	Maturity    *string `json:"maturity"`
}

type StartupDetail struct {
	StartupList
	CreatedAt      *string   `json:"created_at"`
	Description    *string   `json:"description"`
	WebsiteUrl     *string   `json:"website_url"`
	SocialMediaURL *string   `json:"social_media_url"`
	ProjectStatus  *string   `json:"project_status"`
	Needs          *string   `json:"needs"`
	Founders       []Founder `json:"founders" gorm:"foreignKey:StartupUUID;references:UUID"`
}

type StartupCreationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
