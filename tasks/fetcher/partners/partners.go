package partners

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

func setupPartnersQuery(skip uint64, limit uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/partners", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for partners: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("skip", fmt.Sprint(skip))
	queryParams.Set("limit", fmt.Sprint(limit))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getPartners(endpoint *url.URL) ([]legacy.PartnerLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for partners: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var partners []legacy.PartnerLegacy
	_, err = utils.SendRequest("get partners", req, &partners, true)
	if err != nil {
		return nil, err
	}
	return partners, nil
}

func postPartners(partnersLegacy []legacy.PartnerLegacy) []models.Partner {
	var partners []models.Partner
	for _, partnerLegacy := range partnersLegacy {
		partner := models.Partner{
			UUID:            uuid.New().String(),
			ID:              &partnerLegacy.ID,
			Name:            partnerLegacy.Name,
			LegalStatus:     partnerLegacy.LegalStatus,
			Address:         partnerLegacy.Address,
			Email:           partnerLegacy.Email,
			Phone:           partnerLegacy.Phone,
			CreatedAt:       partnerLegacy.CreatedAt,
			Description:     partnerLegacy.Description,
			PartnershipType:    partnerLegacy.PartnershipType,
		}
		var counter int64
		initializers.DB.Table("partners").Where("id=?", *partner.ID).Count(&counter)
		if counter == 0 {
			initializers.DB.Create(&partner)
		}
		partners = append(partners, partner)
	}

	return partners
}

func UpdatePartners(skip uint64, limit uint64) ([]models.Partner, error) {
	var err error = nil

	endpoint, err := setupPartnersQuery(skip, limit)
	if err != nil {
		return nil, err
	}

	partnersLegacy, err := getPartners(endpoint)
	if err != nil {
		return nil, err
	}

	partners := postPartners(partnersLegacy)
	return partners, err
}
