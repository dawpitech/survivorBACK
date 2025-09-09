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
	StartupId   *uint   `json:"startup_id"`
	Description *string `json:"description"`
}

type NewsCreationRequest = struct {
	Title string `json:"title"`
}
