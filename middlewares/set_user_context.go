package middlewares

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUserContext(jwtSecret string) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, _ := stripBearer(context.Request.Header.Get("Authorization"))

		tokenClaims, _ := jwt.ParseWithClaims(
			token,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)

		if tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if ok && tokenClaims.Valid {
				context.Set("user_id", claims.ID)
				context.Set("user_email", claims.Email)

				context.Request = setToContext(context, "user_id", claims.ID)
				context.Request = setToContext(context, "user_email", claims.Email)
			}
		}

		context.Next()
	}
}

func setToContext(c *gin.Context, key interface{}, value interface{}) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
}
