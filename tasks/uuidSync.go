package tasks

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"fmt"
)

func syncNewsUUID(news *models.NewsDetails) {
	if news.StartupUUID != nil {
		return
	}

	var startup models.StartupDetail
	initializers.DB.Where("id=?", news.StartupID).Find(&startup)

	if startup.UUID == "" {
		return
	}
	if rst := initializers.DB.Model(&news).Update("startup_uuid", startup.UUID); rst.Error != nil {
		fmt.Printf("Couldn't update db with re-sync startup UUID on news %s\n", news.UUID)
		return
	}
}

func syncInvestorUUIDs(user *models.User) {
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

func syncFounderUUIDs(user *models.User) {
	if user.FounderID == nil || user.FounderUUID != nil {
		return
	}
	var founderFound models.Founder
	initializers.DB.Where("id=?", user.FounderID).Find(&founderFound)

	if founderFound.UUID == "" {
		return
	}
	if result := initializers.DB.Model(&user).Update("founder_uuid", founderFound.UUID); result.Error != nil {
		fmt.Printf("Couldn't update db with re-sync founder UUID on user %s\n", user.UUID)
		return
	}
}

func SyncUUIDs() {
	var users []models.User
	var news []models.NewsDetails

	if result := initializers.DB.Find(&users); result.Error != nil {
		fmt.Printf("Couldn't fetch users from db to run UUID Sync Task!\n%s\n", result.Error.Error())
		return
	}
	if result := initializers.DB.Find(&news); result.Error != nil {
		fmt.Printf("Couldn't fetch news from db to run UUID Sync Task!\n%s\n", result.Error.Error())
		return
	}

	for _, user := range users {
		syncInvestorUUIDs(&user)
		syncFounderUUIDs(&user)
	}
	for _, singleNews := range news {
		syncNewsUUID(&singleNews)
	}
}
