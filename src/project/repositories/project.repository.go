package repositories

import (
	att "grape/src/attachment/types"
	"grape/src/common/repositories"
	r "grape/src/project/dto/request"
	e "grape/src/project/entities"

	"gorm.io/gorm"
)

type ProjectRelation string

const (
	Organization ProjectRelation = "Organization"
	Attachments  ProjectRelation = "Attachments"
	Owner        ProjectRelation = "Owner"
)

type ProjectRepository struct {
	db *gorm.DB

	repositories.CommonRepository[e.ProjectEntity, r.ProjectDto, ProjectRelation]
}

func (c *ProjectRepository) Build(dto *r.ProjectDto, relations ...ProjectRelation) *gorm.DB {
	tx := c.db.Model(&e.ProjectEntity{})

	var required []ProjectRelation
	c.applyFilter(tx, dto, required)
	c.attachRelations(tx, dto, append(relations, required...))

	return tx
}

func (c *ProjectRepository) applyFilter(tx *gorm.DB, dto *r.ProjectDto, _ []ProjectRelation) {
	if len(dto.ProjectIds) != 0 {
		tx.Where(`projects.uuid IN ?`, dto.ProjectIds)
	}

	if len(dto.Query) != 0 {
		tx.Where(`projects.name ILIKE ?`, "%"+dto.Query+"%")
	}
}

func (c *ProjectRepository) attachRelations(tx *gorm.DB, _ *r.ProjectDto, relations []ProjectRelation) {
	for _, r := range relations {
		switch r {
		case Attachments:
			tx.Joins(`LEFT JOIN attachments ON attachments.attachable_type = ? AND attachments.attachable_id = projects.id`, att.Project)

		default:
			tx.Joins(string(r))
		}
	}
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{}
}
