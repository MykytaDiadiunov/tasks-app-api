package requests

import "go-rest-api/internal/domain"

type CreateProjectRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

func (r CreateProjectRequest) ToDomainModel() (interface{}, error) {
	return domain.Project{
		Title:       r.Title,
		Description: r.Description,
	}, nil
}
