package models

type StartupList struct {
	ID          uint `json:"id"`
	Name        string `json:"name"`
	LegalStatus *string `json:"legal_status"`
	Address     *string	`json:"address"`
	Email       string `json:"email"`
	Phone       *string `json:"phone"`
	Sector      *string `json:"sector"`
	Maturity    *string	`json:"maturity"`
}

type StartupDetail struct {
	StartupList
	CreatedAt      *string `json:"created_at"`
	Description    *string `json:"description"`
	WebsiteUrl     *string `json:"website_url"`
	SocialMediaURL *string `json:"social_media_url"`
	ProjectStatus  *string `json:"project_status"`
	Needs          *string `json:"needs"`
	Founders       []struct {
		Name      string `json:"name"`
		ID        uint `json:"id"`
		StartupID uint `json:"startup_id"`
	} `json:"founders"`
}
