package routes

type GenericUUIDFromPath struct {
	UUID string `path:"uuid" validate:"required"`
}
