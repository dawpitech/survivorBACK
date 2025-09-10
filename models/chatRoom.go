package models

type ChatRoom struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`

	FirstPartyUUID  string `json:"first_party_uuid"`
	SecondPartyUUID string `json:"second_party_uuid"`
}
