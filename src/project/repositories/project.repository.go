package repositories

import (
	att "grape/src/attachment/types"
	"grape/src/common/repositories"
	r "grape/src/project/dto/request"
	e "grape/src/project/entities"
	"grape/src/project/types"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ProjectRelation string

const (
	Organization ProjectRelation = "Organization"
	Attachments  ProjectRelation = "Attachments"
	Owner        ProjectRelation = "Owner"
)

type ProjectRepositoryT = repositories.CommonRepository[e.ProjectEntity, r.ProjectDto, ProjectRelation]

type projectRepository struct {
	db *gorm.DB
}

func (c *projectRepository) Model() *e.ProjectEntity {
	return e.NewProjectEntity()
}

func (c *projectRepository) Build(dto *r.ProjectDto, relations ...ProjectRelation) *gorm.DB {
	tx := c.db.Model(c.Model()).Where(`projects.organization_id = ?`, dto.CurrentUser.Organization.ID)

	var required []ProjectRelation
	c.applyFilter(tx, dto, required)
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *projectRepository) applyFilter(tx *gorm.DB, dto *r.ProjectDto, _ []ProjectRelation) {
	if len(dto.ProjectIds) != 0 {
		tx.Where(`projects.uuid IN ?`, dto.ProjectIds)
	}

	if len(dto.Query) != 0 {
		tx.Where(`projects.name ILIKE ?`, "%"+dto.Query+"%")
	}

	if len(dto.Statuses) != 0 {
		tx.Where(`projects.status IN ?`, lo.Map(dto.Statuses, func(str string, _ int) types.ProjectStatusEnum {
			return types.Active.Value(str)
		}))
	}

	if len(dto.Types) != 0 {
		tx.Where(`projects.type IN ?`, lo.Map(dto.Types, func(str string, _ int) types.ProjectTypeEnum {
			return types.Html.Value(str)
		}))
	}
}

func (c *projectRepository) attachRelations(tx *gorm.DB, _ *r.ProjectDto, relations []ProjectRelation) {
	for _, r := range relations {
		switch r {
		case Attachments:
			tx.Joins(`LEFT JOIN attachments ON attachments.attachable_type = ? AND attachments.attachable_id = projects.id`, att.Project)

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *projectRepository) sortBy(tx *gorm.DB, dto *r.ProjectDto, _ []ProjectRelation) {
	switch dto.SortBy {
	case "name":
	case "order":
	case "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *projectRepository) Create(dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	var order int64
	c.Build(dto).Select(`MAX(projects.order) AS "order"`).Group("projects.id").Scan(&order)

	entity.Order = int(order) + 1
	entity.OwnerID = dto.CurrentUser.ID
	entity.OrganizationID = dto.CurrentUser.Organization.ID

	entity.SetType(body.(*r.ProjectCreateDto).Type)
	return c.db.Model(c.Model()).Create(entity)
}

func (c *projectRepository) Update(dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	options := body.(*r.ProjectUpdateDto)

	entity.SetType(options.Type)
	entity.SetStatus(options.Status)

	return c.db.Model(entity).Updates(entity)
}

func (c *projectRepository) Delete(dto *r.ProjectDto, entity *e.ProjectEntity) *gorm.DB {
	// TODO: Add transaction with recursive delete related entities
	return c.db.Model(c.Model()).Delete(entity)
}

func NewProjectRepository(db *gorm.DB) *ProjectRepositoryT {
	return repositories.NewRepository(&projectRepository{db})
}
