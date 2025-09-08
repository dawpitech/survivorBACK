package routes

type GetStartupRequest = GenericUUIDFromPath

type StartupCreationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateStartupRequest struct {
	UUID           string  `path:"uuid" validate:"required"`
	Name           *string `json:"name"`
	LegalStatus    *string `json:"legal_status"`
	Address        *string `json:"address"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Sector         *string `json:"sector"`
	Maturity       *string `json:"maturity"`
	Description    *string `json:"description"`
	WebsiteUrl     *string `json:"website_url"`
	SocialMediaURL *string `json:"social_media_url"`
	ProjectStatus  *string `json:"project_status"`
	Needs          *string `json:"needs"`
}
