package middlewares

import (
	"github.com/alonelegion/go_graphql_api/controllers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Claims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

func stripBearer(token string) (string, error) {
	if len(token) > 6 && strings.ToLower(token[0:7]) == "bearer" {
		return token[7:], nil
	}
	return token, nil
}

func RequiredLoggedIn(jwtSecret string) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := stripBearer(context.Request.Header.Get("Authorization"))
		if err != nil {
			controllers.HTTPResponse(context, http.StatusUnauthorized, err.Error(), nil)
			context.Abort()
			return
		}

		tokenClaims, err := jwt.ParseWithClaims(
			token,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)
		if err != nil {
			controllers.HTTPResponse(context, http.StatusUnauthorized, err.Error(), nil)
			context.Abort()
			return
		}

		if tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)

			if ok && tokenClaims.Valid {
				context.Set("user_id", claims.ID)
				context.Set("user_email", claims.Email)

				context.Next()
				return
			}
		}

		controllers.HTTPResponse(context, http.StatusUnauthorized, "Unauthorized", nil)
		context.Abort()
		return
	}
}
