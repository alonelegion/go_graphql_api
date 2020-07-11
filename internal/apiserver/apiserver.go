package apiserver

import (
	"fmt"
	"github.com/alonelegion/go_graphql_api/internal/graph"
	"github.com/alonelegion/go_graphql_api/internal/middlewares"
	"log"
	"net/http"

	"github.com/alonelegion/go_graphql_api/configs"
	"github.com/alonelegion/go_graphql_api/internal/controllers"
	"github.com/alonelegion/go_graphql_api/internal/email_client/mailgun_client"
	"github.com/alonelegion/go_graphql_api/internal/general/hmac_hash"
	"github.com/alonelegion/go_graphql_api/internal/general/random_string"

	pwdModel "github.com/alonelegion/go_graphql_api/internal/models/reset_password"
	"github.com/alonelegion/go_graphql_api/internal/models/user"
	"github.com/alonelegion/go_graphql_api/internal/repositories/password_reset"
	"github.com/alonelegion/go_graphql_api/internal/repositories/user_repository"
	"github.com/alonelegion/go_graphql_api/internal/services/auth_service"
	"github.com/alonelegion/go_graphql_api/internal/services/email_service"
	"github.com/alonelegion/go_graphql_api/internal/services/user_service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq" // Postgres setup
)

var (
	router = gin.Default()
)

func Start() {

	// Setup swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup configs
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config := configs.GetConfig()

	// Connect to Database
	db, err := gorm.Open(
		config.Postgres.Dialect(),
		config.Postgres.GetPostgresConnectionInfo(),
	)
	if err != nil {
		panic(err)
	}

	// Migration
	db.AutoMigrate(&user.User{}, &pwdModel.ResetPassword{})
	defer db.Close()

	// Setup EmailClient
	emailClient := mailgun_client.NewMailGunClient(config)

	// Setup Repositories
	userRepo := user_repository.NewUserRepository(db)
	passRepo := password_reset.NewPasswordResetRepository(db)

	// Setup general
	randStr := random_string.NewRandomString()
	hmac := hmac_hash.NewHMAC(config.HMACKey)

	// Setup Services
	userService := user_service.NewUserService(userRepo, passRepo, randStr, hmac, config.Pepper)
	authService := auth_service.NewAuthService(config.JWTSecret)
	emailService := email_service.NewEmailService(emailClient)

	// Setup controllers
	userController := controllers.NewUserController(userService, authService, emailService)

	// Setup middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	router.GET("/ping", func(context *gin.Context) { context.String(http.StatusOK, "pong") })
	router.GET("/graphql", graph.PlayGroundHandler("/query"))
	router.POST("/query",
		middlewares.SetUserContext(config.JWTSecret),
		graph.GraphqlHandler(userService, authService, emailService))

	api := router.Group("/api")

	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)
	api.POST("/forgot_password", userController.ForgotPassword)
	api.POST("/update_password", userController.ResetPassword)

	userGroup := api.Group("/users")

	userGroup.GET("/:id", userController.GetById)

	account := api.Group("/account")
	account.Use(middlewares.RequiredLoggedIn(config.JWTSecret))
	{
		account.GET("/profile", userController.GetProfile)
		account.PUT("/profile", userController.Update)
	}

	// Start
	port := fmt.Sprintf(":%s", config.Port)
	_ = router.Run(port)
}
