package routes

type GetNewsRequest = GenericUUIDFromPath
type DeleteNewsRequest = GenericUUIDFromPath
type GetNewsPictureRequest = GenericUUIDFromPath
type ResetNewsPictureRequest = GenericUUIDFromPath

type NewsUpdateRequest = struct {
	GenericUUIDFromPath
	Title       *string `json:"title"`
	Location    *string `json:"location"`
	Category    *string `json:"category"`
	StartupUUID *string `json:"startup_uuid"`
	Description *string `json:"description"`
}

type NewsCreationRequest = struct {
	Title       string `json:"title"`
	StartupUUID string `json:"startup_uuid"`
}
