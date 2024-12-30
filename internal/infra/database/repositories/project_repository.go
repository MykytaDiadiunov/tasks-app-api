package repositories

import (
	"database/sql"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/logger"
)

type ProjectRepository interface {
	FindById(id uint64) (domain.Project, error)
	FindByCreatorId(creatorId uint64, page, limit int32) (domain.Projects, error)
	Save(project domain.Project) (domain.Project, error)
	Update(project domain.Project) (domain.Project, error)
	Delete(id uint64) error
}

type project struct {
	Id          uint64 `db:"id, omitempty"`
	Title       string `db:"title"`
	Description string `db:"description"`
	CreatorId   uint64 `db:"creator_id"`
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return projectRepository{
		db: db,
	}
}

func (pr projectRepository) FindById(id uint64) (domain.Project, error) {
	sqlCommand := `SELECT * FROM projects WHERE id=$1`
	projectModel := project{}
	err := pr.db.QueryRow(sqlCommand, id).Scan(
		&projectModel.Id,
		&projectModel.Title,
		&projectModel.Description,
		&projectModel.CreatorId,
	)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}

	return pr.modelToDomain(projectModel), nil
}

func (pr projectRepository) FindByCreatorId(creatorId uint64, page, limit int32) (domain.Projects, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	sqlCommand := `SELECT * FROM projects WHERE creator_id=$1 LIMIT $2 OFFSET $3`

	rows, err := pr.db.Query(sqlCommand, creatorId, limit, offset)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Projects{}, err
	}
	defer rows.Close()

	projects := []domain.Project{}
	for rows.Next() {
		projectModel := project{}
		err = rows.Scan(
			&projectModel.Id,
			&projectModel.Title,
			&projectModel.Description,
			&projectModel.CreatorId,
		)
		if err != nil {
			logger.Logger.Error(err)
			return domain.Projects{}, err
		}
		projects = append(projects, pr.modelToDomain(projectModel))
	}

	var total uint64
	totalSqlCommand := `SELECT COUNT(*) FROM projects WHERE creator_id = $1;`
	err = pr.db.QueryRow(totalSqlCommand, creatorId).Scan(&total)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Projects{}, err
	}

	var pages int32
	if total > 0 {
		pages = (int32(total) + limit - 1) / limit
	}

	return domain.Projects{
		Projects:    projects,
		Total:       total,
		CurrentPage: page,
		LastPage:    pages,
	}, nil
}

func (pr projectRepository) Save(project domain.Project) (domain.Project, error) {
	projectModel := pr.domainToModel(project)
	sqlCommand := `INSERT INTO projects (title, description, creator_id) VALUES ($1, $2, $3) RETURNING id`

	err := pr.db.QueryRow(sqlCommand, project.Title, project.Description, project.CreatorId).Scan(&projectModel.Id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}
	return pr.modelToDomain(projectModel), nil
}

func (pr projectRepository) Update(project domain.Project) (domain.Project, error) {
	projectModel := pr.domainToModel(project)
	sqlCommand := `UPDATE projects SET title=$1, description=$2 WHERE id=$3`

	_, err := pr.db.Exec(sqlCommand, projectModel.Title, projectModel.Description, projectModel.Id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}

	updatedProject, err := pr.FindById(projectModel.Id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.Project{}, err
	}
	return updatedProject, nil
}

func (pr projectRepository) Delete(id uint64) error {
	sqlCommand := `DELETE FROM projects WHERE id=$1`

	_, err := pr.db.Exec(sqlCommand, id)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

func (pr projectRepository) modelToDomain(p project) domain.Project {
	return domain.Project{
		Id:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		CreatorId:   p.CreatorId,
	}
}

func (pr projectRepository) domainToModel(p domain.Project) project {
	return project{
		Id:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		CreatorId:   p.CreatorId,
	}
}
