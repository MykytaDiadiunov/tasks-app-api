package resources

import "go-rest-api/internal/domain"

type ProjectDto struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorId   uint64 `json:"creator_id"`
}

type ProjectsDto struct {
	Projects    []ProjectDto `json:"projects"`
	Total       uint64       `json:"total"`
	CurrentPage int32        `json:"current_page"`
	LastPage    int32        `json:"last_page"`
}

func (d ProjectDto) DomainToDto(project domain.Project) ProjectDto {
	return ProjectDto{
		Id:          project.Id,
		Title:       project.Title,
		Description: project.Description,
		CreatorId:   project.CreatorId,
	}
}

func (d ProjectsDto) DomainToDto(projects domain.Projects) ProjectsDto {
	result := make([]ProjectDto, len(projects.Projects))

	for i := range projects.Projects {
		result[i] = ProjectDto{}.DomainToDto(projects.Projects[i])
	}

	return ProjectsDto{
		Projects:    result,
		Total:       projects.Total,
		CurrentPage: projects.CurrentPage,
		LastPage:    projects.LastPage,
	}
}
