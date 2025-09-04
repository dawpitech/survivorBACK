package fetcher

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher/startups"
	"log"
)

func updateStartups() {
	startupList, err := startups.UpdateStartupList(0, 10)
	if err != nil {
		log.Println("Unable to update db for startup: ", err)
	}

	for _, startup := range startupList {
		log.Println("Display startup")
		startupDetail, err := startups.UpdateSingleStartups(uint64(startup.ID))
		if err != nil {
			log.Println("Unable to update db for startup: ", err)
		}

		for _, founder := range startupDetail.Founders {
			_, err = startups.UpdateFounderImage(uint64(startup.ID), uint64(founder.ID))
			if err != nil {
				log.Println("Unable to update db for founder image: ", err)
			}
		}
	}
}

func updateUsers() {
}

func UpdateData() {
	updateStartups()
	updateUsers()
}
