package news

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

func setupNewsQuery(skip uint64, limit uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/news", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for news: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("skip", fmt.Sprint(skip))
	queryParams.Set("limit", fmt.Sprint(limit))

	endpoint.RawQuery = queryParams.Encode()

	return endpoint, nil
}

func getNews(endpoint *url.URL) ([]legacy.NewsLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for news: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var news []legacy.NewsLegacy
	_, err = utils.SendRequest("get news", req, &news, true)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func postNewsList(news []legacy.NewsLegacy) []models.News {
	var newsList []models.News
	for _, newsLegacy := range news {
		news := models.NewsDetails{
			News: models.News{
				UUID:        uuid.New().String(),
				ID:          &newsLegacy.ID,
				Location: newsLegacy.Location,
				Title: newsLegacy.Title,
				Category: newsLegacy.Category,
				StartupId: newsLegacy.StartupId,
			},
		}
		var counter int64
		initializers.DB.Table("news_details").Where("id=?", *news.ID).Count(&counter)
		if counter == 0 {
			initializers.DB.Create(&news)
		}
		newsList = append(newsList, news.News)
	}

	return newsList
}

func UpdateNewsList(skip uint64, limit uint64) ([]models.News, error) {
	var err error = nil

	endpoint, err := setupNewsQuery(skip, limit)
	if err != nil {
		return nil, err
	}

	newsListLegacy, err := getNews(endpoint)
	if err != nil {
		return nil, err
	}

	newsList := postNewsList(newsListLegacy)
	return newsList, err
}
