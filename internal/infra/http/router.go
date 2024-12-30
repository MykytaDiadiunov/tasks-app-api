package http

import (
	"errors"
	"go-rest-api/config/container"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/controllers"
	"go-rest-api/internal/infra/http/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func CreateRouter(con container.Container) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/ping", func(apiRouter chi.Router) {
			apiRouter.Get("/", pingHandler())
			apiRouter.Handle("/*", notFoundJson())
		})
		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, con.SessionController, con.AuthMw)
				})
				apiRouter.Route("/user", func(apiRouter chi.Router) {
					apiRouter.Use(con.AuthMw)
					UserRouter(apiRouter, con)
				})
				apiRouter.Route("/project", func(apiRouter chi.Router) {
					apiRouter.Use(con.AuthMw)
					ProjectRouter(apiRouter, con)
				})
			})
		})
	})

	return router
}

func AuthRouter(r chi.Router, sc controllers.SessionController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			sc.Register(),
		)
		apiRouter.Post(
			"/login",
			sc.Login(),
		)
		apiRouter.With(amw).Delete(
			"/logout",
			sc.Logout(),
		)
	})
}

func UserRouter(r chi.Router, con container.Container) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/me",
			con.UserController.FindMe(),
		)
		apiRouter.Put(
			"/me/update/avatar",
			con.UserController.UpdateUserAvatar(),
		)
		apiRouter.Get(
			"/me/projects",
			con.ProjectController.GetMyProjects(),
		)
		// apiRouter.Get(
		// 	"/email/confirm/{token}",
		// 	uc.ConfirmUserEmailByEmailConfirmationToken(),
		// )
	})
}

func ProjectRouter(r chi.Router, con container.Container) {
	pathObjMw := middlewares.PathObjectMiddleware(con.ProjectService)
	isOwnerMw := middlewares.IsOwnerMiddleware[domain.Project]()
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/{projectId}",
			con.ProjectController.FindProjectById(),
		)
		apiRouter.Post(
			"/",
			con.ProjectController.CreateProject(),
		)
		apiRouter.With(pathObjMw).With(isOwnerMw).Put(
			"/{projectId}",
			con.ProjectController.UpdateProjecTitleAndDescription(),
		)
		apiRouter.With(pathObjMw).With(isOwnerMw).Delete(
			"/{projectId}",
			con.ProjectController.DeleteProjectById(),
		)
	})
}

func notFoundJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.NotFound(w, errors.New("resource Not Found"))
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.Ok(w)
	}
}
