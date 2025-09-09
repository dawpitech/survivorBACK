package models

type News struct {
	UUID     string  `json:"uuid" gorm:"type:uuid;primary_key"`
	ID       *uint   `json:"-" gorm:"unique;index"` // legacy
	Location *string `json:"location"`
	Title    string  `json:"title"`
	Category *string `json:"category"`

	StartupID   *uint   `json:"-"`
	StartupUUID *string `json:"startup_uuid"`
}

type NewsDetails struct {
	News
	Description string `json:"description"`

	NewsPicture *NewsPicture `json:"-" gorm:"foreignKey:NewsUUID;references:UUID"`
}

type NewsPicture struct {
	NewsUUID string `json:"news_uuid" gorm:"type:uuid;primary_key"`
	Picture  []byte `json:"picture" gorm:"type:bytea"`
}
