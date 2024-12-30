package container

import (
	"go-rest-api/config"
	"go-rest-api/internal/app"
	"go-rest-api/internal/infra/database"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/filesystem"
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
	app.ProjectService
}

type Controllers struct {
	controllers.UserController
	controllers.SessionController
	controllers.ProjectController
}

type Middleware struct {
	AuthMw func(http.Handler) http.Handler
}

func New() Container {
	cfg := config.GetConfiguration()
	tknAuth := jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)
	db := database.New(cfg)

	cloudinaryService := filesystem.NewCloudinaryService(cfg)
	// imageService := filesystem.NewImageStorageService("file_storage")

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	projectRepo := repositories.NewProjectRepository(db)

	userService := app.NewUserService(userRepo, cfg, cloudinaryService)
	sessionService := app.NewSessionService(sessionRepo, userService, tknAuth)
	projectService := app.NewProjectService(projectRepo)

	userController := controllers.NewUserController(userService)
	sessionController := controllers.NewSessionController(sessionService, userService)
	projectController := controllers.NewProjectController(projectService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, sessionService, userService)

	return Container{
		Services: Services{
			userService,
			sessionService,
			projectService,
		},
		Controllers: Controllers{
			userController,
			sessionController,
			projectController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}
