package container

import (
	"go-rest-api/config"
	"go-rest-api/internal/app"
	"go-rest-api/internal/infra/database"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/http/controllers"
	"go-rest-api/internal/infra/http/middlewares"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type Container struct {
	Services
	Controllers
	Middleware
}

type Services struct {
	app.UserService
	app.SessionService
}

type Controllers struct {
	controllers.UserController
	controllers.SessionController
}

type Middleware struct {
	AuthMw func(http.Handler) http.Handler
}

func New() Container {
	cfg := config.GetConfiguration()
	tknAuth := jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)
	db := database.New(cfg)

	// cloudinaryService := filesystem.NewCloudinaryService(cfg)
	// imageService := filesystem.NewImageStorageService("file_storage")

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	userService := app.NewUserService(userRepo, cfg)
	sessionService := app.NewSessionService(sessionRepo, userService, tknAuth)

	userController := controllers.NewUserController(userService)
	sessionController := controllers.NewSessionController(sessionService, userService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, sessionService, userService)

	return Container{
		Services: Services{
			userService,
			sessionService,
		},
		Controllers: Controllers{
			userController,
			sessionController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}
