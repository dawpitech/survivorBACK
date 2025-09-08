package investor

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func setupSingleInvestorsImageQuery(investorId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/investors/%d/image", os.Getenv("API_URL"), investorId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for investorImage: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("investor_id", fmt.Sprint(investorId))

	endpoint.RawQuery = queryParams.Encode()
	return endpoint, nil
}

func getSingleInvestorImage(entrypoint *url.URL) (string, error) {
	var investorImage string

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single investor: ", err)
		return investorImage, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	result, err := utils.SendRequest("get for single investor", req, &investorImage, false)
	if err != nil {
		return investorImage, err
	}
	investorImage = result.(string)
	return investorImage, nil
}

func UpdateInvestorImage(investorId uint64) (string, error) {
	var investorImage string

	endpoint, err := setupSingleInvestorsImageQuery(investorId)
	if err != nil {
		return investorImage, err
	}

	investorImage, err = getSingleInvestorImage(endpoint)
	if err != nil {
		return investorImage, err
	}
	return investorImage, nil
}
