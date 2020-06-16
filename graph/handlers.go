package graph

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/alonelegion/go_graphql_api/graph/generated"
	"github.com/alonelegion/go_graphql_api/services/auth_service"
	"github.com/alonelegion/go_graphql_api/services/email_service"
	"github.com/alonelegion/go_graphql_api/services/user_service"
	"github.com/gin-gonic/gin"
)

func GraphqlHandler(
	user user_service.UserService,
	auth auth_service.AuthService,
	es email_service.EmailService) gin.HandlerFunc {
	conf := generated.Config{
		Resolvers: &Resolver{
			UserService:  user,
			AuthService:  auth,
			EmailService: es,
		},
	}
	exec := generated.NewExecutableSchema(conf)
	h := handler.GraphQL(exec)
	return func(context *gin.Context) {
		h.ServeHTTP(context.Writer, context.Request)
	}
}

func PlayGroundHandler(path string) gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", path)
	return func(context *gin.Context) {
		h.ServeHTTP(context.Writer, context.Request)
	}
}
