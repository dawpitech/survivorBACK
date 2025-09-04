package models

type Founder struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"id" gorm:"unique;index"` // legacy

	Name string `json:"name"`

	StartupUUID string         `json:"startup_uuid"`
	StartupID   uint          `json:"startup_id"` // legacy
	Startup     *StartupDetail `json:"startup" gorm:"foreignKey:StartupUUID;references:UUID"`
}
