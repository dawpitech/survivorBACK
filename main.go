package main

import (
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/tasks"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/wI2L/fizz"
	"time"
)

func init() {
	fmt.Println("Initializing...")

	initializers.LoadEnvs(true)
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:80", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	tonic.SetErrorHook(jujerr.ErrHook)
	fizzRouter := fizz.NewFromEngine(router)

	registerRoutes(fizzRouter)

	err = router.Run("0.0.0.0:24680")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
