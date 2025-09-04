package main

import (
	"FranceDeveloppe/JEB-backend/controllers"
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/middlewares"
	"FranceDeveloppe/JEB-backend/tasks"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	router.Use(cors.Default())

	defaultRoutes := router.Group("/api")
	legacyRoutes := defaultRoutes.Group("/legacy/")

	// Auth routes
	defaultRoutes.POST("/auth/signup", controllers.CreateUser)
	defaultRoutes.POST("/auth/login", controllers.LoginUser)

	// User management routes
	defaultRoutes.GET("/users", controllers.GetAllUsers)
	defaultRoutes.PUT("/users", controllers.CreateNewUser)
	defaultRoutes.GET("/user/:uuid", controllers.GetUser)
	defaultRoutes.DELETE("/user/:uuid", controllers.DeleteUser)
	defaultRoutes.PATCH("/user/:uuid", controllers.UpdateUser)

	// TODO/ all legacy Routes will be removed
	legacyRoutes.Use(middlewares.EnsureIncomingFromLocalhost)
	//legacyRoutes.POST("/createUser", controllers.CreateUserFromLegacy)
	//legacyRoutes.POST("/createInvestor", controllers.CreateInvestorFromLegacy)

	err = router.Run("0.0.0.0:24680")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
