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
		fmt.Printf("Couldn't fetch users from db to run UUID Sync Task!\n%s\n", result.Error.Error())
		return
	}

	for _, user := range users {
		if user.UserPicture == nil {
			utils.ResetUserPicture(&user)
		}
	}
}
