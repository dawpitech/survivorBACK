package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"proto/backendAPI/controllers"
	"proto/backendAPI/controllers/legacy"
	"proto/backendAPI/initializers"
	"proto/backendAPI/middlewares"
	"proto/backendAPI/tasks"
)

func init() {
	fmt.Println("Initializing...")

	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("Now launching...")

	tasks.RunTasksInBackground()

	var router = gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defaultRoutes := router.Group("/api")
	legacyRoutes := defaultRoutes.Group("/legacy/")

	defaultRoutes.POST("/auth/signup", controllers.CreateUser)
	defaultRoutes.POST("/auth/login", controllers.LoginUser)
	defaultRoutes.GET("/user/profile", middlewares.CheckAuth, controllers.GetUser)

	legacyRoutes.Use(middlewares.EnsureIncomingFromLocalhost)
	legacyRoutes.POST("/createUser", legacy.CreateUserFromLegacy)
	legacyRoutes.POST("/createInvestor", legacy.CreateInvestorFromLegacy)

	err = router.Run("0.0.0.0:24680")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
