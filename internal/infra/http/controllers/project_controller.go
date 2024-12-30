package controllers

import (
	"errors"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/requests"
	"go-rest-api/internal/infra/http/resources"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProjectController struct {
	projectService app.ProjectService
}

func NewProjectController(projectService app.ProjectService) ProjectController {
	return ProjectController{
		projectService: projectService,
	}
}

func (c ProjectController) FindProjectById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		if projectId == "" {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		numericProjectId, err := strconv.ParseUint(projectId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		project, err := c.projectService.FindById(numericProjectId)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		projectDto := resources.ProjectDto{}
		Success(w, projectDto.DomainToDto(project))
	}
}

func (c ProjectController) GetMyProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		var numericPage, numericLimit int64
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		if page == "" || limit == "" {
			numericPage = 1
			numericLimit = 1
		}

		numericPage, pErr := strconv.ParseInt(page, 10, 32)
		numericLimit, lErr := strconv.ParseInt(limit, 10, 32)
		if pErr != nil || lErr != nil {
			BadRequest(w, errors.New("invalid page or limit"))
			return
		}

		projects, err := c.projectService.FindByCreatorId(user.Id, int32(numericPage), int32(numericLimit))
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Success(w, resources.ProjectsDto{}.DomainToDto(projects))
	}
}

func (c ProjectController) CreateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		projectBody, err := requests.Bind(r, requests.CreateProjectRequest{}, domain.Project{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		projectBody.CreatorId = user.Id

		createdProject, err := c.projectService.Save(projectBody)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		projectDto := resources.ProjectDto{}
		Success(w, projectDto.DomainToDto(createdProject))
	}
}

func (c ProjectController) UpdateProjecTitleAndDescription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		if projectId == "" {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		numericProjectId, err := strconv.ParseUint(projectId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		projectBody, err := requests.Bind(r, requests.CreateProjectRequest{}, domain.Project{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		projectBody.Id = numericProjectId

		updatedProject, err := c.projectService.Update(projectBody)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		projectDto := resources.ProjectDto{}
		Success(w, projectDto.DomainToDto(updatedProject))
	}
}

func (c ProjectController) DeleteProjectById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		if projectId == "" {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		numericProjectId, err := strconv.ParseUint(projectId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid projectId"))
			return
		}

		err = c.projectService.Delete(numericProjectId)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
