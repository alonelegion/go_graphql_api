package user_servece

import (
	"github.com/alonelegion/go_graphql_api/general/hmac_hash"
	rdms "github.com/alonelegion/go_graphql_api/general/random_string"
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
	PassRepo   password_reset.Repository
	RandStr    rdms.RandomString
	hmac       hmac_hash.HMAC
	pepper     string
}

func NewUserService(
	repo user_repository.Repository,
	pwd password_reset.Repository,
	rds rdms.RandomString,
	hmac hmac_hash.HMAC,
	pepper string) UserService {

	return &userService{
		Repository: repo,
		PassRepo:   pwd,
		RandStr:    rds,
		hmac:       hmac,
		pepper:     pepper,
	}
}

func (u userService) GetById(id uint) (*user.User, error) {
	panic("implement me")
}

func (u userService) GetByEmail(email string) (*user.User, error) {
	panic("implement me")
}

func (u userService) Create(u2 *user.User) error {
	panic("implement me")
}

func (u userService) Update(u2 *user.User) error {
	panic("implement me")
}

func (u userService) HashPassword(rawPassword string) (string, error) {
	panic("implement me")
}

func (u userService) ComparePassword(rawPassword string, passwordFromDB string) error {
	panic("implement me")
}

func (u userService) InitiateResetPassword(email string) (string, error) {
	panic("implement me")
}

func (u userService) CompleteUpdatePassword(token, newPassword string) (*user.User, error) {
	panic("implement me")
}
