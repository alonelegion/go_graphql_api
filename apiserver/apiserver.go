package apiserver

import (
	"github.com/alonelegion/go_graphql_api/graph"
	"log"
	"net/http"

	"github.com/alonelegion/go_graphql_api/configs"
	"github.com/alonelegion/go_graphql_api/controllers"
	"github.com/alonelegion/go_graphql_api/email_client/mailgun_client"
	"github.com/alonelegion/go_graphql_api/general/hmac_hash"
	"github.com/alonelegion/go_graphql_api/general/random_string"

	pwdDomain "github.com/alonelegion/go_graphql_api/models/reset_password"
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/alonelegion/go_graphql_api/repositories/password_reset"
	"github.com/alonelegion/go_graphql_api/repositories/user_repository"
	"github.com/alonelegion/go_graphql_api/services/auth_service"
	"github.com/alonelegion/go_graphql_api/services/email_service"
	"github.com/alonelegion/go_graphql_api/services/user_service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	db.AutoMigrate(&user.User{}, &pwdDomain.ResetPassword{})
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
	user_controller := controllers.NewUserController(userService, authService, emailService)

	// Setup middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	router.GET("/ping", func(context *gin.Context) { context.String(http.StatusOK, "pong") })
	router.GET("/graphql", graph.PlayGroundHandler("/query"))
	router.POST("/query")
}
