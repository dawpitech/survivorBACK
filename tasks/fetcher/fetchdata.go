package fetcher

import (
	event "FranceDeveloppe/JEB-backend/tasks/fetcher/events"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/investors"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/news"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/partners"
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
			//TODO: better log policy
			break
		}

		for _, startup := range startupList {
			startupDetail, err := startups.UpdateSingleStartups(uint64(*startup.ID))
			if err != nil {
				//TODO: better log policy
			}

			for _, founder := range startupDetail.Founders {
				startups.UpdateFounderImage(uint64(*startup.ID), uint64(founder.ID))
			}
		}
		startupList, err = startups.UpdateStartupList(uint64(startIndex), uint64(nbToFetch))
	}
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
		for _, singleInvestor := range investorList {
			investor.UpdateInvestorImage(uint64(*singleInvestor.ID))
		}
		investorList, err = investor.UpdateInvestor(uint64(startIndex), uint64(nbToFetch))
	}
}

func updatePartners() {
	startIndex := 0
	nbToFetch := 10
	partnerList, err := partner.UpdatePartners(uint64(startIndex), uint64(nbToFetch))

	for partnerList != nil {
		startIndex += nbToFetch
		if err != nil {
			log.Println("Unable to update db for partner: ", err)
			continue
		}
		partnerList, err = partner.UpdatePartners(uint64(startIndex), uint64(nbToFetch))
	}
}

func updateEvent() {
	startIndex := 0
	nbToFetch := 10
	eventList, err := event.UpdateEvent(uint64(startIndex), uint64(nbToFetch))

	for eventList != nil {
		startIndex += nbToFetch
		if err != nil {
			log.Println("Unable to update db for event: ", err)
			continue
		}
		for _, singleEvent := range eventList {
			event.UpdateEventImage(uint64(*singleEvent.ID))
		}
		eventList, err = event.UpdateEvent(uint64(startIndex), uint64(nbToFetch))
	}
}

func updateUsers() {
	userList, err := user.UpdateUsers()
	if err != nil {
		log.Println("Unable to update db for user: ", err)
		return
	}
	for _, singleUser := range userList {
		user.UpdateUserImage(uint64(singleUser.ID))
	}
}

func updateNews() {
	startIndex := 0
	nbToFetch := 10
	newsList, err := news.UpdateNewsList(uint64(startIndex), uint64(nbToFetch))

	for newsList != nil {
		startIndex += nbToFetch
		if err != nil {
			log.Println("Unable to update db for news: ", err)
			continue
		}
		for _, singleNews := range newsList {
			_, err := news.UpdateSingleNews(uint64(*singleNews.ID))
			if err != nil {
				log.Println("Unable to update db for startup: ", err)
			}
			news.UpdateNewsImage(uint64(*singleNews.ID))

		}
		newsList, err = news.UpdateNewsList(uint64(startIndex), uint64(nbToFetch))
	}
}

func UpdateData() {
	updateStartups()
	updateInvestors()
	updatePartners()
	updateNews()
	updateEvent()
	updateUsers()
	log.Println("DB legacy entirely fetched")
}
