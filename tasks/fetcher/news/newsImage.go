package news

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func setupSingleNewsImageQuery(newsId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/news/%d/image", os.Getenv("API_URL"), newsId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for newsImage: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("news_id", fmt.Sprint(newsId))

	endpoint.RawQuery = queryParams.Encode()
	return endpoint, nil
}

func getSingleNewsImage(entrypoint *url.URL) (string, error) {
	var newsImage string

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single news: ", err)
		return newsImage, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	result, err := utils.SendRequest("get for single news", req, &newsImage, false)
	if err != nil {
		return newsImage, err
	}
	newsImage = result.(string)
	return newsImage, nil
}

func postNewsImage(newsId uint64, image string) error {
	news := models.NewsDetails{}
	result := initializers.DB.Where("id=?", newsId).First(&news)
	if result.Error != nil {
		log.Printf("Error finding news with ID %d: %s", newsId, result.Error)
		return result.Error
	}

	imageBytes := []byte(image)

	newsPicture := models.NewsPicture{
		NewsUUID: news.UUID,
		Picture:  imageBytes,
	}

	var existingPicture models.NewsPicture
	var counter int64
	result = initializers.DB.Table("news_pictures").Where("news_uuid=?", news.UUID).Count(&counter)

	if counter == 0 {
		result = initializers.DB.Create(&newsPicture)
		if result.Error != nil {
			log.Printf("Error creating news picture for news %s: %s", news.UUID, result.Error)
			return result.Error
		}
		return nil
	}

	result = initializers.DB.Model(&existingPicture).Where("news_uuid=?", news.UUID).Update("picture", imageBytes)
	if result.Error != nil {
		log.Printf("Error updating news picture for news %s: %s", news.UUID, result.Error)
		return result.Error
	}

	return nil
}

func UpdateNewsImage(newsId uint64) (string, error) {
	var newsImage string

	endpoint, err := setupSingleNewsImageQuery(newsId)
	if err != nil {
		return newsImage, err
	}

	newsImage, err = getSingleNewsImage(endpoint)
	if err != nil {
		return newsImage, err
	}
	postNewsImage(newsId, newsImage)
	return newsImage, nil
}
