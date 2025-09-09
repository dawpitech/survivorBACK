package utils

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"bytes"
	"fmt"
	"github.com/M1chlCZ/identicon-generator-go"
	"hash/crc32"
	"image/png"
)

func ResetNewsPicture(news *models.NewsDetails) {
	pp := identicon.GenerateIdenticonWithConfig(int(crc32.ChecksumIEEE([]byte(news.UUID))), identicon.Config{
		Width:     512,
		Height:    512,
		GridSize:  8,
		Grayscale: true,
	})

	var buf bytes.Buffer
	if err := png.Encode(&buf, pp); err != nil {
		fmt.Println(err.Error())
		return
	}

	newsPicture := models.NewsPicture{
		NewsUUID: news.UUID,
		Picture:  buf.Bytes(),
	}
	if err := initializers.DB.Save(&newsPicture); err.Error != nil {
		fmt.Println(err.Error)
		return
	}
}

func ResetEventPicture(event *models.Event) {
	pp := identicon.GenerateIdenticonWithConfig(int(crc32.ChecksumIEEE([]byte(event.UUID))), identicon.Config{
		Width:     512,
		Height:    512,
		GridSize:  8,
		Grayscale: true,
	})

	var buf bytes.Buffer
	if err := png.Encode(&buf, pp); err != nil {
		fmt.Println(err.Error())
		return
	}

	eventPicture := models.EventPicture{
		EventUUID: event.UUID,
		Picture:   buf.Bytes(),
	}
	if err := initializers.DB.Save(&eventPicture); err.Error != nil {
		fmt.Println(err.Error)
		return
	}
}

func ResetUserPicture(user *models.User) {
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
	if err := initializers.DB.Save(&userPicture); err.Error != nil {
		fmt.Println(err.Error)
		return
	}
}
