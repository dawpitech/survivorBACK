package models

type Founder struct {
	ID          *uint  `json:"id" gorm:"unique"`
	UUID        string `json:"uuid" gorm:"primary_key"`
	Name        string `json:"name"`
	StartupID   *uint  `json:"startup_id"`
	StartupUUID string `json:"startup_uuid"`
}
