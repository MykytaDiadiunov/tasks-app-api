package app

import (
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/logger"
)

type ProjectService interface {
	FindById(id uint64) (domain.Project, error)
	FindByCreatorId(creatorId uint64, page, limit int32) (domain.Projects, error)
	Save(project domain.Project) (domain.Project, error)
	Update(project domain.Project) (domain.Project, error)
	Delete(id uint64) error
}

type projectService struct {
	projectRepository repositories.ProjectRepository
}

func NewProjectService(projectRepository repositories.ProjectRepository) ProjectService {
	return projectService{
		projectRepository: projectRepository,
	}
}

func (p projectService) FindById(id uint64) (domain.Project, error) {
	project, err := p.projectRepository.FindById(id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}
	return project, err
}

func (p projectService) FindByCreatorId(creatorId uint64, page int32, limit int32) (domain.Projects, error) {
	projects, err := p.projectRepository.FindByCreatorId(creatorId, page, limit)
	if err != nil {
		logger.Logger.Error(err)
		return projects, err
	}
	return projects, nil
}

func (p projectService) Save(project domain.Project) (domain.Project, error) {
	createdProject, err := p.projectRepository.Save(project)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}
	return createdProject, nil
}

func (p projectService) Update(project domain.Project) (domain.Project, error) {
	updatedProject, err := p.projectRepository.Update(project)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}
	return updatedProject, nil
}

func (p projectService) Delete(id uint64) error {
	err := p.projectRepository.Delete(id)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}
