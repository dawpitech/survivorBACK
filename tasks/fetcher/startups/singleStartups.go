package startups

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"FranceDeveloppe/JEB-backend/models/legacy"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
)

func setupSingleStartupQuery(startupId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/startups/%d", os.Getenv("API_URL"), startupId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for single startup: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("startup_id", fmt.Sprint(startupId))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getSingleStartup(entrypoint *url.URL) (legacy.StartupDetailLegacy, error) {
	startup := legacy.StartupDetailLegacy{}

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single startup: ", err)
		return startup, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	_, err = utils.SendRequest("get for single startup", req, &startup, true)
	if err != nil {
		return startup, err
	}
	return startup, nil
}

func UpdateSingleStartups(startupsId uint64) (legacy.StartupDetailLegacy, error) {
	var err error = nil
	startup := legacy.StartupDetailLegacy{}

	endpoint, err := setupSingleStartupQuery(startupsId)
	if err != nil {
		return startup, err
	}

	startup, err = getSingleStartup(endpoint)
	if err != nil {
		return startup, err
	}
	log.Println("Startup info: ", startup)

	return startup, nil
}
