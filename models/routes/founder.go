package routes

type GetFounderRequest = GenericUUIDFromPath

type FounderCreationRequest struct {
	Name string `json:"name" validate:"required"`
}

type FounderUpdateRequest struct {
	GenericUUIDFromPath
	Name        *string `json:"name"`
	StartupUUID *string `json:"startup_uuid"`
}
