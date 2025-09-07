package routes

type GetStartupRequest = GenericUUIDFromPath

type StartupCreationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
