package startups

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/legacy"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
)

func setupStartupsQuery(skip uint64, limit uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/startups", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for startups: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("skip", fmt.Sprint(skip))
	queryParams.Set("limit", fmt.Sprint(limit))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getStartups(endpoint *url.URL) ([]legacy.StartupListLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for startups: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var startups []legacy.StartupListLegacy
	_, err = utils.SendRequest("get startups", req, &startups, true)
	if err != nil {
		return nil, err
	}
	return startups, nil
}

func postStartupList(startups []legacy.StartupListLegacy) error {
	for _, startupLegacy := range(startups) {
		startup := models.StartupList(startupLegacy)
		initializers.DB.Create(startup)
	}

	return nil
}

func UpdateStartup(skip uint64, limit uint64) ([]legacy.StartupListLegacy, error) {
	var err error = nil

	endpoint, err := setupStartupsQuery(skip, limit)
	if err != nil {
		return nil, err
	}

	startupList, err := getStartups(endpoint)
	if err != nil {
		return startupList, err
	}

	err = postStartupList(startupList)
	if err != nil {
		return startupList, err
	}
	log.Println("Startups: ", startupList)
	return startupList, err
}
