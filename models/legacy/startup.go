package legacy

type StartupListLegacy struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	LegalStatus *string `json:"legal_status"`
	Address     *string `json:"address"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone"`
	Sector      *string `json:"sector"`
	Maturity    *string `json:"maturity"`
}

type StartupDetailLegacy struct {
	StartupListLegacy
	CreatedAt      *string         `json:"created_at"`
	Description    *string         `json:"description"`
	WebsiteUrl     *string         `json:"website_url"`
	SocialMediaURL *string         `json:"social_media_url"`
	ProjectStatus  *string         `json:"project_status"`
	Needs          *string         `json:"needs"`
	Founders       []FounderLegacy `json:"founders"`
}
