package tasks

import (
	"fmt"
	"proto/backendAPI/initializers"
	"proto/backendAPI/models"
)

func syncInvestorUUIDs(user models.User) {
	if user.InvestorID == nil || user.InvestorUUID != nil {
		return
	}
	var investorFound models.Investor
	initializers.DB.Where("id=?", user.InvestorID).Find(&investorFound)

	if investorFound.UUID == "" {
		return
	}
	if result := initializers.DB.Model(&user).Update("investor_uuid", investorFound.UUID); result.Error != nil {
		fmt.Printf("Couldn't update db with re-sync investor UUID on user %s\n", user.UUID)
		return
	}
}

func SyncUUIDs() {
	var users []models.User

	if result := initializers.DB.Find(&users); result.Error != nil {
		fmt.Println("Couldn't fetch users from db to run UUID Sync Task!")
		return
	}

	for _, user := range users {
		syncInvestorUUIDs(user)
	}
}
