package user_service

import (
	"errors"
	"fmt"
	"github.com/alonelegion/go_graphql_api/general/hmac_hash"
	rdms "github.com/alonelegion/go_graphql_api/general/random_string"
	pwd "github.com/alonelegion/go_graphql_api/models/reset_password"
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/alonelegion/go_graphql_api/repositories/password_reset"
	"github.com/alonelegion/go_graphql_api/repositories/user_repository"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (u *userService) GetById(id uint) (*user.User, error) {
	if id == 0 {
		return nil, errors.New("id param is required")
	}
	user, err := u.Repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) GetByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, errors.New("email (string) is required")
	}
	user, err := u.Repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) Create(user *user.User) error {
	hashedPass, err := u.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPass
	return u.Repository.Create(user)
}

func (u *userService) Update(user *user.User) error {
	return u.Repository.Update(user)
}

func (u *userService) HashPassword(rawPassword string) (string, error) {
	passAndPepper := rawPassword + u.pepper
	hashed, err := bcrypt.GenerateFromPassword([]byte(passAndPepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}

func (u *userService) ComparePassword(rawPassword string, passwordFromDB string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(passwordFromDB),
		[]byte(rawPassword+u.pepper),
	)
}

// Функция присваения токена пользователю для обновления пароля
func (u *userService) InitiateResetPassword(email string) (string, error) {
	user, err := u.Repository.GetByEmail(email)
	if err != nil {
		return "", err
	}

	token, err := u.RandStr.GenerateToken()
	if err != nil {
		return "", err
	}

	hashedToken := u.hmac.Hash(token)

	pwd := pwd.ResetPassword{
		UserID: user.ID,
		Token:  hashedToken,
	}

	err = u.PassRepo.Create(&pwd)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) CompleteUpdatePassword(token, newPassword string) (*user.User, error) {
	hashedToken := u.hmac.Hash(token)

	pwr, err := u.PassRepo.GetOneByToken(hashedToken)
	if err != nil {
		return nil, err
	}

	if time.Now().Sub(pwr.CreatedAt) > (1 * time.Hour) {
		return nil, errors.New("invalid Token")
	}

	user, err := u.Repository.GetByID(pwr.UserID)
	if err != nil {
		return nil, err
	}

	hashedPass, err := u.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPass
	if err = u.Repository.Update(user); err != nil {
		return nil, err
	}

	if err = u.PassRepo.Delete(pwr.ID); err != nil {
		fmt.Println("Не удалось удалить запись сброса пароля", pwr.ID, err.Error())
	}
	return user, nil
}
