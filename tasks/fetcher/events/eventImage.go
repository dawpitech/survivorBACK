package event

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func setupSingleEventsImageQuery(eventId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/events/%d/image", os.Getenv("API_URL"), eventId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for eventImage: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("event_id", fmt.Sprint(eventId))

	endpoint.RawQuery = queryParams.Encode()
	return endpoint, nil
}

func getSingleEventImage(entrypoint *url.URL) (string, error) {
	var eventImage string

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single event: ", err)
		return eventImage, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	result, err := utils.SendRequest("get for single event", req, &eventImage, false)
	if err != nil {
		return eventImage, err
	}
	eventImage = result.(string)
	return eventImage, nil
}

func postEventImage(eventId uint64, image string) error {
	event := models.Event{}
	result := initializers.DB.Where("id=?", eventId).First(&event)
	if result.Error != nil {
		log.Printf("Error finding event with ID %d: %s", eventId, result.Error)
		return result.Error
	}

	imageBytes := []byte(image)

	eventPicture := models.EventPicture{
		EventUUID: event.UUID,
		Picture:  imageBytes,
	}

	var existingPicture models.EventPicture
	var counter int64
	result = initializers.DB.Table("event_pictures").Where("event_uuid=?", event.UUID).Count(&counter)

	if counter == 0 {
		result = initializers.DB.Create(&eventPicture)
		if result.Error != nil {
			log.Printf("Error creating event picture for event %s: %s", event.UUID, result.Error)
			return result.Error
		}
		return nil
	}

	result = initializers.DB.Model(&existingPicture).Where("event_uuid=?", event.UUID).Update("picture", imageBytes)
	if result.Error != nil {
		log.Printf("Error updating event picture for event %s: %s", event.UUID, result.Error)
		return result.Error
	}

	return nil
}

func UpdateEventImage(eventId uint64) (string, error) {
	var eventImage string

	endpoint, err := setupSingleEventsImageQuery(eventId)
	if err != nil {
		return eventImage, err
	}

	eventImage, err = getSingleEventImage(endpoint)
	if err != nil {
		return eventImage, err
	}
	postEventImage(eventId, eventImage)
	return eventImage, nil
}

