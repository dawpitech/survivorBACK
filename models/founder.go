package models

type Founder struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"-" gorm:"unique;index"` // legacy

	Name string `json:"name"`

	StartupUUID string         `json:"startup_uuid"`
	StartupID   uint           `json:"-"` // legacy
	Startup     *StartupDetail `json:"startup,omitempty" gorm:"foreignKey:StartupUUID;references:UUID"`

	FounderPicture *FounderPicture `json:"-" gorm:"foreignKey:FounderUUID;references:UUID"`
}

type FounderPicture struct {
	FounderUUID string `json:"founder_uuid" gorm:"type:uuid;primary_key"`
	Picture     []byte `json:"picture" gorm:"type:bytea"`
}
