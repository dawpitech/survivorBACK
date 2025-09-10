package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"time"
)

func CreateChatRoom(_ *gin.Context, in *routes.CreateRoomRequest) (*models.ChatRoom, error) {
	var chatRoomSO models.ChatRoom
	initializers.DB.
		Where("first_party_uuid=?", in.FirstPartyUUID).
		Where("second_party_uuid=?", in.SecondPartyUUID).
		Find(&chatRoomSO)

	var chatRoomRO models.ChatRoom
	initializers.DB.
		Where("first_party_uuid=?", in.SecondPartyUUID).
		Where("second_party_uuid=?", in.FirstPartyUUID).
		Find(&chatRoomRO)

	if chatRoomSO.UUID != "" || chatRoomRO.UUID != "" {
		return nil, errors.NewNotValid(nil, "Room already existing with those users")
	}

	chatRoom := models.ChatRoom{
		UUID:            uuid.New().String(),
		FirstPartyUUID:  in.FirstPartyUUID,
		SecondPartyUUID: in.SecondPartyUUID,
	}

	if err := initializers.DB.Create(&chatRoom); err.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &chatRoom, nil
}

func SendMessageInChatRoom(_ *gin.Context, in *routes.CreateMessageRequest) (*models.ChatMessage, error) {
	var chatRoom models.ChatRoom
	initializers.DB.Where("uuid=?", in.UUID).Find(&chatRoom)

	if chatRoom.UUID == "" {
		return nil, errors.NewNotFound(nil, "No room with given UUID")
	}

	chatMessage := models.ChatMessage{
		UUID:         uuid.New().String(),
		ChatRoomUUID: in.UUID,
		SenderUUID:   in.SenderUUID,
		ReceiverUUID: in.ReceiverUUID,
		Content:      in.Content,
		SentAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := initializers.DB.Save(&chatMessage); err.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &chatMessage, nil
}

func GetRoomMessages(_ *gin.Context, in *routes.GetRoomMessagesRequest) (*[]models.ChatMessage, error) {
	var chatRoom models.ChatRoom
	initializers.DB.Where("uuid=?", in.UUID).Find(&chatRoom)

	if chatRoom.UUID == "" {
		return nil, errors.NewNotFound(nil, "No room with given UUID")
	}

	var chatRoomMessages []models.ChatMessage

	if err := initializers.DB.Where("chat_room_uuid=?", in.UUID).Find(&chatRoomMessages); err.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &chatRoomMessages, nil
}

func GetAllChatRooms(_ *gin.Context, _ *struct{}) (*[]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom

	if err := initializers.DB.Find(&chatRooms); err.Error != nil {
		return nil, errors.New("Internal server error")
	}
	return &chatRooms, nil
}
