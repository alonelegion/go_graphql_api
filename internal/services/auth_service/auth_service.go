package auth_service

import (
	"github.com/alonelegion/go_graphql_api/internal/models/user"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

type Claims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

type authService struct {
	jwtSecret string
}

type AuthService interface {
	IssueToken(u user.User) (string, error)
	ParseToken(token string) (*Claims, error)
}

func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

func (a *authService) IssueToken(u user.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // 24 часа

	claims := Claims{
		u.Email,
		u.ID,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Go GraphQL API",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return tokenClaims.SignedString([]byte(a.jwtSecret))
}

func (a *authService) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwtSecret), nil
		},
	)
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
