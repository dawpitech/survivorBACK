package startups

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
)

func setupFounderImageQuery(startupId uint64, founderId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/startups/%d/founders/%d/image", os.Getenv("API_URL"), startupId, founderId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for founder image: ", err.Error())
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("startup_id", fmt.Sprint(startupId))
	queryParams.Set("founder_id", fmt.Sprint(founderId))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getFounderImage(entrypoint *url.URL) (string, error) {
	var founderImage string

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for founder image: ", err.Error())
		return founderImage, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	result, err := utils.SendRequest("founder image", req, &founderImage, false)
	if err != nil {
		return founderImage, err
	}
	founderImage = result.(string)

	return founderImage, nil
}

func postFounderImage(founderId uint64, image string) error {
	founder := models.Founder{}
	result := initializers.DB.Where("id=?", founderId).First(&founder)
	if result.Error != nil {
		log.Printf("Error finding founder with ID %d: %s", founderId, result.Error)
		return result.Error
	}

	imageBytes := []byte(image)

	founderPicture := models.FounderPicture{
		FounderUUID: founder.UUID,
		Picture:     imageBytes,
	}

	var existingPicture models.FounderPicture
	var counter int64
	result = initializers.DB.Table("founder_pictures").Where("founder_uuid=?", founder.UUID).Count(&counter)

	if counter == 0 {
		result = initializers.DB.Create(&founderPicture)
		if result.Error != nil {
			log.Printf("Error creating founder picture for founder %s: %s", founder.UUID, result.Error)
			return result.Error
		}
		return nil
	}

	result = initializers.DB.Model(&existingPicture).Where("founder_uuid=?", founder.UUID).Update("picture", imageBytes)
	if result.Error != nil {
		log.Printf("Error updating founder picture for founder %s: %s", founder.UUID, result.Error)
		return result.Error
	}

	return nil
}

func UpdateFounderImage(startupsId uint64, founderId uint64) (string, error) {
	var err error = nil
	var founderImage string

	endpoint, err := setupFounderImageQuery(startupsId, founderId)
	if err != nil {
		return founderImage, err
	}

	founderImage, err = getFounderImage(endpoint)
	if err != nil {
		return founderImage, err
	}

	postFounderImage(founderId, founderImage)
	return founderImage, nil
}
