package legacy

type StartupList struct {
	ID          uint
	Name        string
	LegalStatus *string
	Address     *string
	Email       string
	Phone       *string
	Sector      *string
	Maturity    *string
}

type StartupDetail struct {
	StartupList
	CreatedAt      *string
	Description    *string
	WebsiteUrl     *string
	SocialMediaURL *string
	ProjectStatus  *string
	Needs          *string
	Founders       []struct {
		Name      string
		ID        uint
		StartupID uint
	}
}
