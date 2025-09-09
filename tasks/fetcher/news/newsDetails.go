package news

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
)

func setupSingleNewsQuery(newsId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/news/%d", os.Getenv("API_URL"), newsId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for single news: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("news_id", fmt.Sprint(newsId))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getSingleNews(entrypoint *url.URL) (legacy.NewsDetailsLegacy, error) {
	news := legacy.NewsDetailsLegacy{}

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single news: ", err)
		return news, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	_, err = utils.SendRequest("get for single news", req, &news, true)
	if err != nil {
		return news, err
	}
	return news, nil
}

func postNewsDetail(newsLegacy legacy.NewsDetailsLegacy) error {

	news := models.NewsDetails{
		News: models.News{
			ID:          &newsLegacy.ID,
			Location:    newsLegacy.Location,
			Title:       newsLegacy.Title,
			Category:    newsLegacy.Category,
			StartupID:   newsLegacy.StartupId,
			StartupUUID: nil,
		},
		Description: newsLegacy.Description,
	}
	var newsUUID string
	initializers.DB.Table("news_details").Select("uuid").Where("id=?", *news.ID).Limit(1).Find(&newsUUID)
	news.UUID = newsUUID

	initializers.DB.Where("id=?", *news.ID).Save(&news)
	return nil
}

func UpdateSingleNews(newsId uint64) (legacy.NewsDetailsLegacy, error) {
	var err error = nil
	news := legacy.NewsDetailsLegacy{}

	endpoint, err := setupSingleNewsQuery(newsId)
	if err != nil {
		return news, err
	}

	news, err = getSingleNews(endpoint)
	if err != nil {
		return news, err
	}

	err = postNewsDetail(news)
	if err != nil {
		return news, err
	}

	return news, nil
}
