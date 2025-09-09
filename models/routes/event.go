package routes

type GetEventRequest = GenericUUIDFromPath
type GetEventPictureRequest = GenericUUIDFromPath
type DeleteEventRequest = GenericUUIDFromPath
type ResetEventPictureRequest = GenericUUIDFromPath

type UpdateEventRequest struct {
	GenericUUIDFromPath
	Name           *string `json:"name"`
	Date           *string `json:"date"`
	Location       *string `json:"location"`
	Description    *string `json:"description"`
	EventType      *string `json:"event_type"`
	TargetAudience *string `json:"target_audience"`
}

type EventCreationRequest struct {
	Name string `json:"name" binding:"required"`
}
