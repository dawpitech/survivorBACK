package tasks

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"bytes"
	"fmt"
	"github.com/M1chlCZ/identicon-generator-go"
	"hash/crc32"
	"image/png"
)

func generatePP(user *models.User) {
	pp := identicon.GenerateIdenticonWithConfig(int(crc32.ChecksumIEEE([]byte(user.UUID))), identicon.Config{
		Width:     512,
		Height:    512,
		GridSize:  5,
		Grayscale: false,
	})

	var buf bytes.Buffer
	if err := png.Encode(&buf, pp); err != nil {
		fmt.Println(err.Error())
		return
	}

	userPicture := models.UserPicture{
		UserUUID: user.UUID,
		Picture:  buf.Bytes(),
	}
	if err := initializers.DB.Create(&userPicture); err.Error != nil {
		fmt.Println(err.Error)
		return
	}
}

func UpdateUsersWithoutPP() {
	var users []models.User

	if result := initializers.DB.Preload("UserPicture").Find(&users); result.Error != nil {
		fmt.Printf("Couldn't fetch users from db to run UUID Sync Task!\n%s\n", result.Error.Error())
		return
	}

	for _, user := range users {
		if user.UserPicture == nil {
			generatePP(&user)
		}
	}
}
