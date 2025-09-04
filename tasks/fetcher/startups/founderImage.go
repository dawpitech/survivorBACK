package startups

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
)

func setupFounderImageQuery(startupId uint64, founderId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/startups/%d/founders/%d/image", os.Getenv("API_URL"), startupId, founderId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for founder image: ", err)
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
		log.Println("Error in get request for founder image: ", err)
		return founderImage, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	_, err = utils.SendRequest("founder image", req, &founderImage, false)
	if err != nil {
		return founderImage, err
	}

	return founderImage, nil
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

	return founderImage, nil
}
