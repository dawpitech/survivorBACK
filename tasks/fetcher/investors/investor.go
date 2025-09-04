package investor

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

func setupInvestorsQuery(skip uint64, limit uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/investors", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for investors: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("skip", fmt.Sprint(skip))
	queryParams.Set("limit", fmt.Sprint(limit))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getInvestors(endpoint *url.URL) ([]legacy.InvestorLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for investors: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var investors []legacy.InvestorLegacy
	_, err = utils.SendRequest("get investors", req, &investors, true)
	if err != nil {
		return nil, err
	}
	return investors, nil
}

func postInvestor(investorsLegacy []legacy.InvestorLegacy) []models.Investor {
	var investors []models.Investor
	for _, investorLegacy := range investorsLegacy {
		investor := models.Investor{
			UUID:            uuid.New().String(),
			ID:              &investorLegacy.ID,
			Name:            investorLegacy.Name,
			LegalStatus:     investorLegacy.LegalStatus,
			Address:         investorLegacy.Address,
			Email:           investorLegacy.Email,
			Phone:           investorLegacy.Phone,
			CreatedAt:       investorLegacy.CreatedAt,
			Description:     investorLegacy.Description,
			InvestorType:    investorLegacy.InvestorType,
			InvestmentFocus: investorLegacy.InvestmentFocus,
		}
		var counter int64
		initializers.DB.Table("investors").Where("id=?", *investor.ID).Count(&counter)
		if counter == 0 {
			initializers.DB.Create(&investor)
		}
		investors = append(investors, investor)
	}

	return investors
}

func UpdateInvestor(skip uint64, limit uint64) ([]models.Investor, error) {
	var err error = nil

	endpoint, err := setupInvestorsQuery(skip, limit)
	if err != nil {
		return nil, err
	}

	investorsLegacy, err := getInvestors(endpoint)
	if err != nil {
		return nil, err
	}

	investors := postInvestor(investorsLegacy)
	return investors, err
}
