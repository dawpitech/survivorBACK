package legacy

type NewsLegacy struct {
	ID        uint    `json:"id"`
	Location  *string `json:"location"`
	Title     string  `json:"title"`
	Category  *string `json:"category"`
	StartupId *uint   `json:"startup_id"`
}

type NewsDetailsLegacy struct {
	NewsLegacy
	Description string `json:"description"`
}
