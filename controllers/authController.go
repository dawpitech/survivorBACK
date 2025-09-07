package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func LoginUser(_ *gin.Context, in *routes.AuthInput) (*routes.AuthResponse, error) {
	var userFound models.User
	initializers.DB.Where("email=?", in.Email).Find(&userFound)

	if userFound.UUID == "" || userFound.Password == nil {
		return nil, errors.NewNotFound(nil, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*userFound.Password), []byte(in.Password)); err != nil {
		return nil, errors.NewNotFound(nil, "Invalid username or password")
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": userFound.UUID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return nil, errors.New("Internal server error")
	}

	response := routes.AuthResponse{
		Token: token,
	}
	return &response, nil
}

func CreateUser(_ *gin.Context, in *routes.AuthInput) (*struct{}, error) {
	var userFound models.User
	if findResult := initializers.DB.Where("email=?", in.Email).Find(&userFound); findResult.Error != nil {
		return nil, errors.New("Internal server error")
	}

	if userFound.Password != nil && userFound.UUID != "" {
		return nil, errors.NewNotValid(nil, "Email already used")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Internal server error")
	}

	password := string(passwordHash)
	if userFound.Password == nil && userFound.UUID != "" {
		if updateResult := initializers.DB.Model(&userFound).Update("password", password); updateResult.Error != nil {
			return nil, errors.New("Internal server error")
		}
		var empty struct{}
		return &empty, nil
	} else {
		// Temporarily disable creation of new account
		return nil, errors.NewForbidden(nil, "Account creation is disabled.")
		/*
			user := models.User{
				UUID:     uuid.New().String(),
				Email:    authInput.Email,
				Password: &password,
			}
			if createResult := initializers.DB.Create(&user); createResult.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"user": user.GetPublicUser()})
		*/
	}
}
