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
	NewsPicture *NewsPicture `json:"news_picture" gorm:"foreignKey:NewsUUID;references:UUID"`
}

type NewsPicture struct {
	NewsUUID string `json:"news_uuid" gorm:"type:uuid;primary_key"`
	Picture  []byte `json:"picture" gorm:"type:bytea"`
}
