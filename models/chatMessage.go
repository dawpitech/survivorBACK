package models

type ChatMessage struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`

	ChatRoomUUID string `json:"chat_room_uuid"`
	SenderUUID   string `json:"sender_uuid"`
	ReceiverUUID string `json:"receiver_uuid"`
	Content      string `json:"content"`
	SentAt       string `json:"sent_at"`
}
