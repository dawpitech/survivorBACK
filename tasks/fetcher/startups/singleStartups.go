package startups

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/legacy"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"

	"github.com/google/uuid"
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

func postStartupDetail(startupLegacy legacy.StartupDetailLegacy) error {
	var founders []models.Founder
	for _, founderLegacy := range startupLegacy.Founders {
		founder := models.Founder{
			UUID:      uuid.New().String(),
			ID:        &founderLegacy.ID,
			Name:      founderLegacy.Name,
			StartupID: &founderLegacy.StartupID,
		}
		founders = append(founders, founder)
	}

	startup := models.StartupDetail{
		StartupList: models.StartupList{
			ID:          &startupLegacy.ID,
			Name:        startupLegacy.Name,
			LegalStatus: startupLegacy.LegalStatus,
			Address:     startupLegacy.Address,
			Email:       startupLegacy.Email,
			Phone:       startupLegacy.Phone,
			Sector:      startupLegacy.Sector,
			Maturity:    startupLegacy.Maturity,
		},
		CreatedAt:      startupLegacy.CreatedAt,
		Description:    startupLegacy.Description,
		WebsiteUrl:     startupLegacy.WebsiteUrl,
		SocialMediaURL: startupLegacy.SocialMediaURL,
		ProjectStatus:  startupLegacy.ProjectStatus,
		Needs:          startupLegacy.Needs,
		Founders:       founders,
	}
	var startupUUID string
	initializers.DB.Table("startup_details").Select("uuid").Where("id=?", *startup.ID).Limit(1).Find(&startupUUID)
	startup.UUID = startupUUID

	for i := range startup.Founders {
		startup.Founders[i].StartupUUID = startupUUID
	}

	for _, founder := range startup.Founders {
		var existingFounder models.Founder
		result := initializers.DB.Where("id = ?", *founder.ID).First(&existingFounder)
		if result.Error != nil {
			initializers.DB.Create(&founder)
		}
	}
	initializers.DB.Where("id=?", *startup.ID).Omit("Founders").Save(&startup)
	return nil
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

	err = postStartupDetail(startup)
	if err != nil {
		return startup, err
	}

	return startup, nil
}
