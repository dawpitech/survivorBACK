package models

type Event struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"id" gorm:"unique;index"` // legacy

	Name           string  `json:"name"`
	Date           *string `json:"date"`
	Location       *string `json:"location"`
	Description    *string `json:"description"`
	EventType      *string `json:"event_type"`
	TargetAudience *string `json:"target_audience"`

	EventPicture *EventPicture `json:"event_picture" gorm:"foreignKey:EventUUID;references:UUID"`
}

type EventPicture struct {
	EventUUID string `json:"event_uuid" gorm:"type:uuid;primary_key"`
	Picture   []byte `json:"picture" gorm:"type:bytea"`
}
