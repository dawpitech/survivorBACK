package main

import (
	"fmt"
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

	if os.Getenv("GIN_MODE") == "debug" {
		fmt.Println("Warning: running migration in 'debug' mode, dev credentials will be available.")
		initializers.DB.Create(&models.User{
			UUID:         "99999999-9999-9999-9999-999999999999",
			ID:           nil,
			Name:         "Dev Local Admin",
			Email:        "dev@francedeveloppe.fr",
			Password:     "$2a$10$5igFsJ1MJXyvgUI42oZWxOyv1ukLssH70t/ig21Bs.D5mPd0gDXtC", // FranceDeveloppe
			Role:         "admin",
			FounderID:    nil,
			FounderUUID:  nil,
			InvestorID:   nil,
			InvestorUUID: nil,
		})
	}

}
