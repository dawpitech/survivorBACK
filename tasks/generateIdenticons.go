package tasks

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/utils"
	"fmt"
)

func UpdateUsersWithoutPP() {
	var users []models.User

	if result := initializers.DB.Preload("UserPicture").Find(&users); result.Error != nil {
		fmt.Printf("Couldn't fetch users from db to run identicons check!\n%s\n", result.Error.Error())
		return
	}

	for _, user := range users {
		if user.UserPicture == nil {
			utils.ResetUserPicture(&user)
		}
	}
}

func UpdateEventsWithoutP() {
	var events []models.Event

	if result := initializers.DB.Preload("EventPicture").Find(&events); result.Error != nil {
		fmt.Printf("Couldn't fetch events from db to run identicons check!\n%s\n", result.Error.Error())
		return
	}

	for _, event := range events {
		if event.EventPicture == nil {
			utils.ResetEventPicture(&event)
		}
	}
}
