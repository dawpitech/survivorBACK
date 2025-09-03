package controllers

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func LoginUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.UUID == "" || userFound.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": userFound.UUID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	if findResult := initializers.DB.Where("email=?", authInput.Email).Find(&userFound); findResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if userFound.Password != nil && userFound.UUID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	password := string(passwordHash)
	if userFound.Password == nil && userFound.UUID != "" {
		if updateResult := initializers.DB.Model(&userFound).Update("password", password); updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userFound.GetPublicUser()})
	} else {
		// Temporarily disable creation of new account
		c.JSON(http.StatusForbidden, gin.H{"error": "Account creation is disabled."})
		return
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
