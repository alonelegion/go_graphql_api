package controllers

import (
	"github.com/alonelegion/go_graphql_api/models/user"
	"github.com/alonelegion/go_graphql_api/services/auth_service"
	"github.com/alonelegion/go_graphql_api/services/email_service"
	"github.com/alonelegion/go_graphql_api/services/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Структура представляет формат запроса входа или регистрации
type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Структура представляет собой возврат пользователя
type UserOutput struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
}

// Структура представляет собой формат запроса на обновление профиля
type UserUpdateInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserController interface {
	Register(*gin.Context)
	Login(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	Update(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

type userController struct {
	user user_service.UserService
	auth auth_service.AuthService
	es   email_service.EmailService
}

func NewUserController(
	user user_service.UserService,
	auth auth_service.AuthService,
	es email_service.EmailService) UserController {
	return &userController{
		user: user,
		auth: auth,
		es:   es,
	}
}

func (ctl *userController) Register(context *gin.Context) {
	// Считать ввод пользователя
	var userInput UserInput
	if err := context.ShouldBindJSON(&userInput); err != nil {
		HTTPResponse(context, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	u := ctl.inputToUser(userInput)

	// Create user
	if err := ctl.user.Create(&u); err != nil {
		HTTPResponse(context, http.StatusInternalServerError, err.Error(), nil)
	}

	// Отправка приветственного сообщения
	if err := ctl.es.Welcome(u.Email); err != nil {
		HTTPResponse(context, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Login
	err := ctl.login(context, &u)
	if err != nil {
		HTTPResponse(context, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

func (ctl *userController) Login(c *gin.Context) {
	// Считать ввод пользователя
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Получить пользователя из базы данных
	user, err := ctl.user.GetByEmail(userInput.Email)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Проверка пароля
	err = ctl.user.ComparePassword(userInput.Password, user.Password)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Логин
	err = ctl.login(c, user)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

func (ctl *userController) GetById(c *gin.Context) {
	panic("implement me")
}

func (ctl *userController) GetProfile(c *gin.Context) {
	panic("implement me")
}

func (ctl *userController) Update(c *gin.Context) {
	panic("implement me")
}

func (ctl *userController) ForgotPassword(c *gin.Context) {
	panic("implement me")
}

func (ctl *userController) ResetPassword(c *gin.Context) {
	panic("implement me")
}

/*
	======== Private Methods =======
*/

func (ctl *userController) inputToUser(input UserInput) user.User {
	return user.User{
		Email:    input.Email,
		Password: input.Password,
	}
}

func (ctl *userController) login(context *gin.Context, u *user.User) error {
	token, err := ctl.auth.IssueToken(*u)
	if err != nil {
		return err
	}

	userOutput := ctl.mapToUserOutput(u)
	out := gin.H{"token": token, "user": userOutput}
	HTTPResponse(context, http.StatusOK, "ok", out)
	return nil
}

func (ctl *userController) mapToUserOutput(u *user.User) *UserOutput {
	return &UserOutput{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Active:    u.Active,
	}
}
