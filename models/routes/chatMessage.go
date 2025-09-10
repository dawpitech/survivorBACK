package routes

type GetRoomMessagesRequest = GenericUUIDFromPath

type CreateMessageRequest struct {
	GenericUUIDFromPath
	SenderUUID   string `json:"sender_uuid"`
	ReceiverUUID string `json:"receiver_uuid"`
	Content      string `json:"content"`
}
