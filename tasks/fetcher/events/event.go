package event

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/legacy"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
)

func setupEventsQuery(skip uint64, limit uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/events", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for events: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("skip", fmt.Sprint(skip))
	queryParams.Set("limit", fmt.Sprint(limit))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getEvents(endpoint *url.URL) ([]legacy.EventLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for events: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var events []legacy.EventLegacy
	_, err = utils.SendRequest("get events", req, &events, true)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func postEvent(eventsLegacy []legacy.EventLegacy) []models.Event {
	var events []models.Event
	for _, eventLegacy := range eventsLegacy {
		event := models.Event{
			UUID:            uuid.New().String(),
			ID:              &eventLegacy.ID,
			Name:            eventLegacy.Name,
			Location: eventLegacy.Location,
			Description:     eventLegacy.Description,
			EventType:    eventLegacy.EventType,
			TargetAudience: eventLegacy.TargetAudience,
		}
		var counter int64
		initializers.DB.Table("events").Where("id=?", *event.ID).Count(&counter)
		if counter == 0 {
			initializers.DB.Create(&event)
		}
		events = append(events, event)
	}

	return events
}

func UpdateEvent(skip uint64, limit uint64) ([]models.Event, error) {
	var err error = nil

	endpoint, err := setupEventsQuery(skip, limit)
	if err != nil {
		return nil, err
	}

	eventsLegacy, err := getEvents(endpoint)
	if err != nil {
		return nil, err
	}

	events := postEvent(eventsLegacy)
	return events, err
}

