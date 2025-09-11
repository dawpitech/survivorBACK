package main

import (
	"FranceDeveloppe/JEB-backend/controllers"
	"FranceDeveloppe/JEB-backend/middlewares"
	"FranceDeveloppe/JEB-backend/models/routes"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

func registerRoutes(fizzRouter *fizz.Fizz) {
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
	usersRoutes := globalRoutes.Group("/users", "User", "This group contains all users endpoints")
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
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
	usersRoutes.GET(
		"founderUUID/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the user with the corresponding Founder UUID"),
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
		tonic.Handler(controllers.GetUserByFounderUUID, 200),
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or user on itself"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UpdateUser, 200),
	)
	usersRoutes.GET(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Get the user's profile picture"),
		},
		tonic.Handler(controllers.GetUserPicture, 200),
	)
	usersRoutes.PUT(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Update the user's profile picture"),
			fizz.InputModel(routes.GenericUUIDFromPath{}),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or user on itself"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UpdateUserPicture, 200),
	)
	usersRoutes.DELETE(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Reset the user's profile picture"),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or user on itself"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.ResetUserPicture, 200),
	)

	// Startup management routes
	startupRoutes := globalRoutes.Group("/startups", "Startup", "This group contains all startups endpoints")
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
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
	startupRoutes.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the startup with the corresponding UUID"),
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.DeleteStartup, 200),
	)
	startupRoutes.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the startup with the corresponding UUID"),
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
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or founders of the startup"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UpdateStartup, 200),
	)
	startupRoutes.PATCH(
		"/:uuid/upViewCount",
		[]fizz.OperationOption{
			fizz.Summary("Up the count of startup views (matching the given UUID) by one"),
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
		tonic.Handler(controllers.AddViewToStartup, 200),
	)
	startupRoutes.GET(
		"/:uuid/file",
		[]fizz.OperationOption{
			fizz.Summary("Get the pitch deck of the startup with the given UUID"),
			fizz.Response("404",
				"File not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetStartupFile, 200),
	)
	startupRoutes.PUT(
		"/:uuid/file",
		[]fizz.OperationOption{
			fizz.Summary("Upload the pitch deck of the startup with the given UUID"),
			fizz.Response("400",
				"Bad request",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or founders of the startup"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UploadStartupFile, 200),
	)

	// Founders management routes
	founderRoutes := globalRoutes.Group("/founders", "Founders", "This group contains all founders endpoints")
	founderRoutes.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all founders"),
		},
		tonic.Handler(controllers.GetAllFounders, 200),
	)
	founderRoutes.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Create a new founder"),
		},
		tonic.Handler(controllers.CreateNewFounder, 200),
	)
	founderRoutes.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the founder with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Founder not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetFounder, 200),
	)
	founderRoutes.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the founder with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Founder not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.DeleteFounder, 200),
	)
	founderRoutes.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the founder with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Founder not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.UpdateFounder, 200),
	)

	// Investors management routes
	investorRoutes := globalRoutes.Group("/investors", "Investors", "This group contains all investors endpoints")
	investorRoutes.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all investors"),
		},
		tonic.Handler(controllers.GetAllInvestors, 200),
	)
	investorRoutes.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Create a new investor"),
		},
		tonic.Handler(controllers.CreateNewInvestor, 200),
	)
	investorRoutes.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the investor with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Investor not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetInvestor, 200),
	)
	investorRoutes.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the investor with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Investor not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.DeleteInvestor, 200),
	)
	investorRoutes.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the investor with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Investor not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.UpdateInvestor, 200),
	)

	// Events management routes
	eventsRoute := globalRoutes.Group("/events", "Events", "This group contains all events endpoints")
	eventsRoute.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all events"),
		},
		tonic.Handler(controllers.GetAllEvents, 200),
	)
	eventsRoute.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Register a new event"),
			fizz.Response(
				"400",
				"Email already used",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.CreateNewEvent, 200),
	)
	eventsRoute.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the event with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Event not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetEvent, 200),
	)
	eventsRoute.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the event with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Event not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.DeleteEvent, 200),
	)
	eventsRoute.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the event with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"Event not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.UpdateEvent, 200),
	)
	eventsRoute.GET(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Get the event picture"),
		},
		tonic.Handler(controllers.GetEventPicture, 200),
	)
	eventsRoute.PUT(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Update the event picture"),
			fizz.InputModel(routes.GenericUUIDFromPath{}),
		},
		tonic.Handler(controllers.UpdateEventPicture, 200),
	)
	eventsRoute.DELETE(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Reset the event picture"),
		},
		tonic.Handler(controllers.ResetEventPicture, 200),
	)

	// News management routes
	newsRoute := globalRoutes.Group("/news", "News", "This group contains all news endpoints")
	newsRoute.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get list of all news"),
		},
		tonic.Handler(controllers.GetAllNews, 200),
	)
	newsRoute.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Register a new news"),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or founders"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.CreateNewNews, 200),
	)
	newsRoute.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get the news with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"News not found",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetNews, 200),
	)
	newsRoute.DELETE(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Delete the news with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"News not found",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.DeleteNews, 200),
	)
	newsRoute.PATCH(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Update the news with the corresponding UUID"),
			fizz.Response(
				"400",
				"Invalid UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Response(
				"404",
				"News not found",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or founders"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UpdateNews, 200),
	)
	newsRoute.GET(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Get the news picture"),
		},
		tonic.Handler(controllers.GetNewsPicture, 200),
	)
	newsRoute.PUT(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Update the news picture"),
			fizz.InputModel(routes.GenericUUIDFromPath{}),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or founders"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.UpdateNewsPicture, 200),
	)
	newsRoute.DELETE(
		"/:uuid/picture",
		[]fizz.OperationOption{
			fizz.Summary("Reset the news picture"),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.ResetNewsPicture, 200),
	)

	// Chat management routes
	chatRoute := globalRoutes.Group("/rooms", "Messages", "This group contains all messages endpoints")
	chatRoute.GET(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Get all chat rooms"),
			fizz.Response(
				"400",
				"Room already existing",
				routes.ErrorOutput{},
				nil,
				nil),
		},
		tonic.Handler(controllers.GetAllChatRooms, 200),
	)
	chatRoute.POST(
		"/",
		[]fizz.OperationOption{
			fizz.Summary("Create a new chat room with the two user based on their user UUID"),
			fizz.Response(
				"400",
				"No room with the given UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or users included in the request"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.CreateChatRoom, 200),
	)
	chatRoute.GET(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Get all message in the correspond room based on the UUID"),
			fizz.Response(
				"400",
				"No room with the given UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or users included in the chat room"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.GetRoomMessages, 200),
	)
	chatRoute.PUT(
		"/:uuid",
		[]fizz.OperationOption{
			fizz.Summary("Create a new message in the correspond room based on the UUID"),
			fizz.Response(
				"400",
				"No room with the given UUID",
				routes.ErrorOutput{},
				nil,
				nil),
			fizz.Security(&openapi.SecurityRequirement{
				"bearerAuth": []string{},
			}),
			fizz.Description("Operation restricted to admins or users included in the chat room"),
		},
		middlewares.CheckAuth,
		tonic.Handler(controllers.SendMessageInChatRoom, 200),
	)

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
}
