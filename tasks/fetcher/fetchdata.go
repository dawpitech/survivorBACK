package fetcher

import (
	investor "FranceDeveloppe/JEB-backend/tasks/fetcher/investors"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/startups"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/users"
	"log"
)

func updateStartups() {
	startIndex := 0
	nbToFetch := 10
	startupList, err := startups.UpdateStartupList(uint64(startIndex), uint64(nbToFetch))

	for startupList != nil {
		startIndex += nbToFetch
		if err != nil {
			log.Println("Unable to update db for startup: ", err)
			break
		}

		for _, startup := range startupList {
			startupDetail, err := startups.UpdateSingleStartups(uint64(*startup.ID))
			if err != nil {
				log.Println("Unable to update db for startup: ", err)
			}

			for _, founder := range startupDetail.Founders {
				startups.UpdateFounderImage(uint64(*startup.ID), uint64(founder.ID))
			}
		}
		startupList, err = startups.UpdateStartupList(uint64(startIndex), uint64(nbToFetch))
	}
	log.Println("DB legacy entirely fetched")
}

func updateInvestors() {
	startIndex := 0
	nbToFetch := 10
	investorList, err := investor.UpdateInvestor(uint64(startIndex), uint64(nbToFetch))

	for investorList != nil {
		startIndex += nbToFetch
		if err != nil {
			log.Println("Unable to update db for investors: ", err)
			continue
		}
		investorList, err = investor.UpdateInvestor(uint64(startIndex), uint64(nbToFetch))
	}
}

func updateUsers() {
	_, err := user.UpdateUsers()
	if err != nil {
		log.Println("Unable to update db for user: ", err)
		return
	}
	// for _, singleUser := range(userList) {
	// 	user.UpdateUserImage(uint64(singleUser.ID))
	// }
}

func UpdateData() {
	updateStartups()
	updateInvestors()
	updateUsers()
}
