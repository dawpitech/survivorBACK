package legacy

type EventLegacy struct {
	ID uint `json:"id"` // legacy

	Name           string  `json:"name"`
	Date           *string `json:"date"`
	Location       *string `json:"location"`
	Description    *string `json:"description"`
	EventType      *string `json:"event_type"`
	TargetAudience *string `json:"target_audience"`
}
