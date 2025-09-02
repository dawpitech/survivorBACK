package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("FranceDeveloppe"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	passwordHashString := string(passwordHash)

	if os.Getenv("GIN_MODE") == "debug" {
		fmt.Println("Warning: running migration in 'debug' mode, dev credentials will be available.")
		initializers.DB.Create(&models.User{
			UUID:         "99999999-9999-9999-9999-999999999999",
			ID:           nil,
			Name:         "Dev Local Admin",
			Email:        "dev@francedeveloppe.fr",
			Password:     &passwordHashString, // FranceDeveloppe
			Role:         "admin",
			FounderID:    nil,
			FounderUUID:  nil,
			InvestorID:   nil,
			InvestorUUID: nil,
		})
	}

}
