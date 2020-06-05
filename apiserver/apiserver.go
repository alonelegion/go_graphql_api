package apiserver

import (
	"github.com/alonelegion/go_graphql_api/configs"
	pwdDomain "github.com/alonelegion/go_graphql_api/models/reset_password"
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
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

}
