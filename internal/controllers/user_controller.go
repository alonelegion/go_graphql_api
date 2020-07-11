package controllers

import (
	"github.com/alonelegion/go_graphql_api/internal/models/user"
	"github.com/alonelegion/go_graphql_api/internal/services/auth_service"
	"github.com/alonelegion/go_graphql_api/internal/services/email_service"
	"github.com/alonelegion/go_graphql_api/internal/services/user_service"
	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"strconv"
	"strings"
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
	us user_service.UserService,
	auth auth_service.AuthService,
	es email_service.EmailService) UserController {
	return &userController{
		user: us,
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
	id, err := ctl.getUserId(c.Param(("id")))
	if err != nil {
		HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	user, err := ctl.user.GetByID(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPResponse(c, http.StatusOK, "ok", userOutput)
}

func (ctl *userController) GetProfile(c *gin.Context) {
	id, exists := c.Get("user_id")
	if exists == false {
		HTTPResponse(c, http.StatusBadRequest, "Неверный идентификатор пользователя", nil)
		return
	}

	user, err := ctl.user.GetByID(id.(uint))
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPResponse(c, http.StatusOK, "ok", userOutput)
}

func (ctl *userController) Update(c *gin.Context) {
	// Получить идентификатор пользователя из контекста
	id, exist := c.Get("user_id")
	if exist == false {
		HTTPResponse(c, http.StatusBadRequest, "Неверный идентификатор пользователя", nil)
		return
	}

	user, err := ctl.user.GetByID(id.(uint))
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Считать ввод пользователя
	var userInput UserUpdateInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Проверка пользователя
	if user.ID != id {
		HTTPResponse(c, http.StatusUnauthorized, "Неравторизованный", nil)
		return
	}

	// Обновление записи пользователя
	user.FirstName = userInput.FirstName
	user.LastName = userInput.LastName
	user.Email = userInput.Email
	if err := ctl.user.Update(user); err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Ответ
	userOutput := ctl.mapToUserOutput(user)
	HTTPResponse(c, http.StatusOK, "ok", userOutput)
}

func (ctl *userController) ForgotPassword(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Присвоение токена пользователю для обновления пароля
	token, err := ctl.user.InitiateResetPassword(input.Email)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Отправка email с токеном для обновления пароля
	if err := ctl.es.ResetPassword(input.Email, token); err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	HTTPResponse(c, http.StatusOK, "Email отправлен", nil)
	return
}

func (ctl *userController) ResetPassword(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token := c.Request.URL.Query().Get("token")
	if token == "" {
		HTTPResponse(c, http.StatusNotFound, "Требуется токен", nil)
		return
	}

	user, err := ctl.user.CompleteUpdatePassword(token, input.Password)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = ctl.login(c, user)
	if err != nil {
		HTTPResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
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

func (ctl *userController) getUserId(userIdParam string) (uint, error) {
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		return 0, errors.New("идентификатор пользователя должен быть числом")
	}
	return uint(userId), nil
}
