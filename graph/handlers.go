package graph

import (
	"github.com/alonelegion/go_graphql_api/services/auth_service"
	"github.com/alonelegion/go_graphql_api/services/email_service"
	"github.com/alonelegion/go_graphql_api/services/user_service"
	"github.com/gin-gonic/gin"
)

func GraphqlHandler(
	user user_service.UserService,
	auth auth_service.AuthService,
	email email_service.EmailService) gin.HandlerFunc {

}
