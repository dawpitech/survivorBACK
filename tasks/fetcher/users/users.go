package user

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

func setupUsersQuery() (*url.URL, error) {
	baseUrl := fmt.Sprintf("%s/users", os.Getenv("API_URL"))
	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Println("Error when parsing base url for users: ", err)
		return nil, err
	}
	return endpoint, nil
}

func getUsers(endpoint *url.URL) ([]legacy.UserLegacy, error) {
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Println("Error in get request for users: ", err)
		return nil, err
	}
	req.Header.Add("X-Group-Authorization", os.Getenv("API_KEY"))

	var users []legacy.UserLegacy
	_, err = utils.SendRequest("get users", req, &users, true)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func postUser(userList []legacy.UserLegacy) error {
	for _, userLegacy := range userList {
		user := models.User{
			UUID:       uuid.New().String(),
			ID:         &userLegacy.ID,
			Email:      userLegacy.Email,
			Name:       userLegacy.Name,
			Role:       userLegacy.Role,
			FounderID:  userLegacy.FounderID,
			InvestorID: userLegacy.InvestorID,
		}
		var counter int64
		initializers.DB.Table("users").Where("id=?", *user.ID).Count(&counter)
		if counter == 0 {
			initializers.DB.Create(&user)
		}
	}
	return nil
}

func UpdateUsers() ([]legacy.UserLegacy, error) {
	endpoint, err := setupUsersQuery()
	if err != nil {
		return nil, err
	}

	users, err := getUsers(endpoint)
	if err != nil {
		return users, err
	}

	err = postUser(users)
	if err != nil {
		return users, err
	}
	return users, err
}
