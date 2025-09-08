package routes

type UserCreationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

type GetUserRequest = GenericUUIDFromPath
type GetUserPictureRequest = GenericUUIDFromPath

type UpdateUserRequest struct {
	UUID     string `path:"uuid" validate:"required"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteUserRequest = GetUserRequest
