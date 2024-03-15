package repositories

import (
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
	Links        ProjectRelation = "Links"
	Owner        ProjectRelation = "Owner"
)

type ProjectRepositoryT = repositories.CommonRepository[*e.ProjectEntity, *r.ProjectDto, ProjectRelation]

type projectRepository struct {
	db *gorm.DB
}

func (c *projectRepository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return c.db
}

func (c *projectRepository) Model() *e.ProjectEntity {
	return e.NewProjectEntity()
}

func (c *projectRepository) Transaction(fc func(*gorm.DB) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		return fc(tx.Model(c.Model()))
	})
}

func (c *projectRepository) Build(db *gorm.DB, dto *r.ProjectDto, relations ...ProjectRelation) *gorm.DB {
	tx := c.conn(db).Model(c.Model()).Where(`projects.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []ProjectRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *projectRepository) applyFilter(tx *gorm.DB, dto *r.ProjectDto, relations []ProjectRelation) []ProjectRelation {
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

	return relations
}

func (c *projectRepository) attachRelations(tx *gorm.DB, _ *r.ProjectDto, relations []ProjectRelation) {
	for _, r := range relations {
		switch r {
		case Links:
		case Attachments:
			tx.Preload(string(r))

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

func (c *projectRepository) Create(tx *gorm.DB, dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	var order int64
	c.Build(tx, dto).Select(`MAX(projects.order) AS "order"`).Group("projects.id").Scan(&order)

	entity.Order = int(order) + 1
	entity.Owner = *dto.CurrentUser
	entity.Organization = dto.CurrentUser.Organization

	entity.SetType(body.(*r.ProjectCreateDto).Type)
	return c.conn(tx).Create(entity)
}

func (c *projectRepository) Update(tx *gorm.DB, dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	options := body.(*r.ProjectUpdateDto)

	entity.SetType(options.Type)
	entity.SetStatus(options.Status)

	return c.conn(tx).Model(entity).Updates(entity)
}

func (c *projectRepository) Delete(tx *gorm.DB, dto *r.ProjectDto, entity *e.ProjectEntity) *gorm.DB {
	// TODO: Add transaction with recursive delete related entities
	return c.conn(tx).Model(c.Model()).Delete(entity)
}

func NewProjectRepository(db *gorm.DB) *ProjectRepositoryT {
	return repositories.NewRepository(&projectRepository{db})
}
