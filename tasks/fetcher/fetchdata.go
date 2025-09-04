package fetcher

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher/startups"
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
			continue
		}

		for _, startup := range startupList {
			startupDetail, err := startups.UpdateSingleStartups(uint64(*startup.ID))
			if err != nil {
				log.Println("Unable to update db for startup: ", err)
			}

			for _, founder := range startupDetail.Founders {
				_, err = startups.UpdateFounderImage(uint64(*startup.ID), uint64(founder.ID))
				if err != nil {
					log.Println("Unable to update db for founder image: ", err)
				}
			}
		}
		startupList, err = startups.UpdateStartupList(uint64(startIndex), uint64(nbToFetch))
	}
	log.Println("DB legacy entirely fetched")
}

func updateUsers() {
}

func UpdateData() {
	updateStartups()
	updateUsers()
}
