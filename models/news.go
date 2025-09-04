package models

type News struct {
	UUID      string  `json:"uuid" gorm:"type:uuid;primary_key"`
	ID        *uint   `json:"id" gorm:"unique;index"` // legacy
	Location  *string `json:"location"`
	Title     string  `json:"title"`
	Category  *string `json:"category"`
	StartupId *uint   `json:"startup_id"`
}

type NewsDetails struct {
	News
	Description string `json:"description"`
}
