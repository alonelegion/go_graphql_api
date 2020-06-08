package user_servece

import (
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/alonelegion/go_graphql_api/repositories/password_reset"
	"github.com/alonelegion/go_graphql_api/repositories/user_repository"
)

type UserService interface {
	GetById(id uint) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Create(*user.User) error
	Update(*user.User) error
	HashPassword(rawPassword string) (string, error)
	ComparePassword(rawPassword string, passwordFromDB string) error
	InitiateResetPassword(email string) (string, error)
	CompleteUpdatePassword(token, newPassword string) (*user.User, error)
}

type userService struct {
	Repository user_repository.Repository
	passRepo   password_reset.Repository
}
