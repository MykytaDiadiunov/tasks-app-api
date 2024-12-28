package controllers

import (
	"errors"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/requests"
	"go-rest-api/internal/infra/http/resources"
	"go-rest-api/internal/infra/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(userService app.UserService) UserController {
	return UserController{userService: userService}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		Success(w, resources.UserDto{}.ToDomainModel(user))
	}
}

func (c UserController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		user, err = c.userService.Save(user)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Created(w, user)
	}
}

func (c UserController) ConfirmUserEmailByEmailConfirmationToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		confirmToken := chi.URLParam(r, "token")

		if confirmToken == "" {
			err := errors.New("invalid confirm token")
			logger.Logger.Error(err)
			BadRequest(w, err)
		}

		user.EmailConfirmationToken = confirmToken

		err := c.userService.ConfirmUserEmail(user)
		if err != nil {
			logger.Logger.Error(err)
			BadRequest(w, err)
		}
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		err := c.userService.Delete(user.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		Ok(w)
	}
}
