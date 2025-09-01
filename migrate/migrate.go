package main

import (
	"fmt"
	"proto/backendAPI/initializers"
	"proto/backendAPI/models"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
