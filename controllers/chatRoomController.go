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

func CreateChatRoom(c *gin.Context, in *routes.CreateRoomRequest) (*models.ChatRoom, error) {
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

	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	if authUser.Role != "admin" && authUser.UUID != in.FirstPartyUUID && authUser.UUID != in.SecondPartyUUID {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

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

func SendMessageInChatRoom(c *gin.Context, in *routes.CreateMessageRequest) (*models.ChatMessage, error) {
	var chatRoom models.ChatRoom
	initializers.DB.Where("uuid=?", in.UUID).Find(&chatRoom)

	if chatRoom.UUID == "" {
		return nil, errors.NewNotFound(nil, "No room with given UUID")
	}

	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	if authUser.Role != "admin" && authUser.UUID != chatRoom.FirstPartyUUID && authUser.UUID != chatRoom.SecondPartyUUID {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

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

func GetRoomMessages(c *gin.Context, in *routes.GetRoomMessagesRequest) (*[]models.ChatMessage, error) {
	var chatRoom models.ChatRoom
	initializers.DB.Where("uuid=?", in.UUID).Find(&chatRoom)

	if chatRoom.UUID == "" {
		return nil, errors.NewNotFound(nil, "No room with given UUID")
	}

	// START AUTH CHECK SECTION
	userInterface, exist := c.Get("currentUser")

	if !exist {
		return nil, errors.New("Internal server error")
	}

	var authUser models.User
	switch u := userInterface.(type) {
	case models.User:
		authUser = u
	case *models.User:
		authUser = *u
	default:
		return nil, errors.New("Internal server error")
	}

	if authUser.Role != "admin" && authUser.UUID != chatRoom.FirstPartyUUID && authUser.UUID != chatRoom.SecondPartyUUID {
		return nil, errors.NewForbidden(nil, "Access Forbidden")
	}
	// END AUTH CHECK SECTION

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
