package user

import (
	"FranceDeveloppe/JEB-backend/tasks/fetcher/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func setupSingleUsersQuery(userId uint64) (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/users/%d/image", os.Getenv("API_URL"), userId)
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for userImage: ", err)
		return nil, err
	}
	queryParams := url.Values{}
	queryParams.Set("user_id", fmt.Sprint(userId))

	endpoint.RawQuery = queryParams.Encode()
	return endpoint, nil
}

func getSingleUser(entrypoint *url.URL) (string, error) {
	var userImage string

	req, err := http.NewRequest("GET", entrypoint.String(), nil)

	if err != nil {
		log.Println("Error in get request for single user: ", err)
		return userImage, nil
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	_, err = utils.SendRequest("get for single user", req, &userImage, false)
	if err != nil {
		return userImage, err
	}
	return userImage, nil
}

// I musn't forget to implement this function to post image
func postUserImage(userId uint64, image string) error {
	return nil
}

func UpdateUserImage(userId uint64) (string, error) {
	var userImage string

	endpoint, err := setupSingleUsersQuery(userId)
	if err != nil {
		return userImage, err
	}

	userImage, err = getSingleUser(endpoint)
	if err != nil {
		return userImage, err
	}
	err = postUserImage(userId, userImage)
	return userImage, nil
}
