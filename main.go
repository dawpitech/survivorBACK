package main

import (
	"FranceDeveloppe/JEB-backend/controllers"
	"FranceDeveloppe/JEB-backend/initializers"
	"FranceDeveloppe/JEB-backend/middlewares"
	"FranceDeveloppe/JEB-backend/models/routes"
	"FranceDeveloppe/JEB-backend/tasks"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

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
		AllowOrigins:     []string{"http://localhost", "http://localhost:80", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	tonic.SetErrorHook(jujerr.ErrHook)
	fizzRouter := fizz.NewFromEngine(router)
	globalRoutes := fizzRouter.Group("/api", "Others", "")

	// Auth routes
	authRoutes := globalRoutes.Group("/auth", "Authentification", "This group contains all sessions authentifications endpoints")
	authRoutes.POST(
		"signup",
		[]fizz.OperationOption{
			fizz.Summary("Sign-up a new user"),
			fizz.Response(
				"400",
				"Email already used",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.CreateUser, 200),
	)
	authRoutes.POST(
		"login",
		[]fizz.OperationOption{
			fizz.Summary("Sign-in as an user"),
			fizz.Response(
				"404",
				"Invalid username or password",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.LoginUser, 200),
	)

	// User management routes
	usersRoutes := globalRoutes.Group("/users", "User access", "This group contains all users endpoints")
	usersRoutes.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all registered users"),
		},
		tonic.Handler(controllers.GetAllUsers, 200),
	)
	usersRoutes.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Register a new user"),
			fizz.Response(
				"400",
				"Email already used",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"User not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.CreateNewUser, 200),
	)
	usersRoutes.GET(
		"/me",
		[]fizz.OperationOption{
			fizz.Summary("Get informations about the logged-in user"),
			fizz.Response(
				"401",
				"Unauthorized",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.GetMe, 200),
	)
	usersRoutes.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the user with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"User not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetUser, 200),
	)
	usersRoutes.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the user with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"User not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.DeleteUser, 200),
	)
	usersRoutes.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the user with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"User not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.UpdateUser, 200),
	)
	usersRoutes.GET(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Get user's profile picture"),
		},
		tonic.Handler(controllers.GetUserPicture, 200),
	)

	// Startup management routes
	startupRoutes := globalRoutes.Group("/startups", "Startup access", "This group contains all startups endpoints")
	startupRoutes.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all startups"),
		},
		tonic.Handler(controllers.GetAllStartups, 200),
	)
	startupRoutes.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Create a new startup"),
			fizz.Response(
				"400",
				"Email already used",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.CreateNewStartup, 200),
	)
	startupRoutes.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the startup with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Startup not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetStartup, 200),
	)
	startupRoutes.DELETE("/:uuid", nil, controllers.DeleteStartup)
	startupRoutes.PATCH("/:uuid", nil, controllers.UpdateStartup)

	fizzRouter.Generator().SetSecuritySchemes(map[string]*openapi.SecuritySchemeOrRef{
		"bearerAuth": {
			SecurityScheme: &openapi.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	})
	infos := &openapi.Info{
		Title:       "JEB Incubator internal API",
		Description: "Internal API used by the JEB incubator platform",
		Version:     "2.1.0",
	}
	fizzRouter.GET("/api/openapi.json", nil, fizzRouter.OpenAPI(infos, "json"))
	fizzRouter.GET("/favicon.ico", nil, func(c *gin.Context) {
		c.File("./favicon.ico")
	})

	err = router.Run("0.0.0.0:24680")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
