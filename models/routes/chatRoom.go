package routes

type CreateRoomRequest struct {
	FirstPartyUUID  string `json:"first_party_uuid" validate:"required"`
	SecondPartyUUID string `json:"second_party_uuid" validate:"required"`
}
