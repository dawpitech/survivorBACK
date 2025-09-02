package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"proto/backendAPI/controllers"
	"proto/backendAPI/controllers/legacy"
	"proto/backendAPI/initializers"
	"proto/backendAPI/middlewares"
)

type user struct {
	UserID     int32 `binding:"required"`
	FounderID  int32
	InvestorID int32
	Name       string
	Email      string
}

func init() {
	fmt.Println("Initializing...")

	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("Now launching...")
	var router = gin.Default()

	router.POST("/auth/signup", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)

	router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUser)

	systemRoutes := router.Group("/legacy/")

	systemRoutes.Use(middlewares.EnsureIncomingFromLocalhost)
	systemRoutes.POST("/createUser", legacy.CreateUserFromLegacy)

	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
