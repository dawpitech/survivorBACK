package fetcher

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher/startups"
	"log"

)

func UpdateData() {
	_, err := startups.UpdateStartup(0, 10)
	if err != nil {
		log.Println("Unable to update db for startup: ", err)
	}
	_, err = startups.UpdateSingleStartups( 8)
	if err != nil {
		log.Println("Unable to update db for startup: ", err)
	}

	_, err = startups.UpdateFounderImage( 8, 13)
	if err != nil {
		log.Println("Unable to update db for founder image: ", err)
	}
}
