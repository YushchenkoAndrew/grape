package repositories

import (
	"fmt"
	"grape/src/common/repositories"
	t "grape/src/common/types"
	r "grape/src/project/dto/request"
	e "grape/src/project/entities"
	"grape/src/project/types"
	st "grape/src/statistic/entities"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ProjectRelation string

const (
	Organization ProjectRelation = "Organization"
	Thumbnail    ProjectRelation = "Thumbnail"
	Attachments  ProjectRelation = "Attachments"
	Redirect     ProjectRelation = "Redirect"
	Links        ProjectRelation = "Links"
	Tags         ProjectRelation = "Tags"
	Owner        ProjectRelation = "Owner"
	Statistic    ProjectRelation = "Statistic"
)

type ProjectRepositoryT = repositories.CommonRepository[*e.ProjectEntity, *r.ProjectDto, ProjectRelation]

type projectRepository struct {
}

func (c *projectRepository) Model() *e.ProjectEntity {
	return e.NewProjectEntity()
}

func (c *projectRepository) Build(db *gorm.DB, dto *r.ProjectDto, relations ...ProjectRelation) *gorm.DB {
	tx := db.Model(c.Model()).Where(`projects.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []ProjectRelation{Thumbnail, Redirect})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *projectRepository) applyFilter(tx *gorm.DB, dto *r.ProjectDto, relations []ProjectRelation) []ProjectRelation {
	if len(dto.ProjectIds) != 0 {
		tx.Where(`projects.uuid IN ?`, dto.ProjectIds)
	}

	if len(dto.Query) != 0 {
		query := "%" + dto.Query + "%"
		tx.Where(
			tx.Session(&gorm.Session{NewDB: true}).Where(`projects.name ILIKE ?`, query).
				Or(`(SELECT COUNT(t.id) FROM tags t WHERE t.taggable_id = projects.id AND t.taggable_type = ? AND t.name ILIKE ?) > 0`, c.Model().TableName(), query),
		)
	}

	if len(dto.Statuses) != 0 {
		tx.Where(`projects.status IN ?`, lo.Map(dto.Statuses, func(str string, _ int) t.StatusEnum {
			return t.Active.Value(str)
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
		case Thumbnail:
			tx.Preload(string(r), "name ILIKE ?", "thumbnail%")

		case Redirect:
			tx.Preload(string(r), "name ILIKE ?", "redirect")

		case Attachments, Links:
			tx.Preload(string(r), func(db *gorm.DB) *gorm.DB { return db.Order(fmt.Sprintf("%s.order ASC", string(r))) })

		case Tags:
			tx.Preload(string(r))

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *projectRepository) sortBy(tx *gorm.DB, dto *r.ProjectDto, _ []ProjectRelation) {
	switch dto.SortBy {
	case "":
		return

	case "name", "order", "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *projectRepository) Create(db *gorm.DB, dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	var order int64
	dto.SortBy = ""
	c.Build(db, dto).Select(`COALESCE(MAX(projects.order), 0) AS "order"`).Scan(&order)

	entity.Order = int(order) + 1
	entity.Owner = dto.CurrentUser
	entity.Organization = &dto.CurrentUser.Organization
	entity.Statistic = st.NewStatisticEntity()

	entity.SetType(body.(*r.ProjectCreateDto).Type)
	return db.Create(entity)
}

func (c *projectRepository) Update(db *gorm.DB, dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	options := body.(*r.ProjectUpdateDto)

	// entity.SetType(options.Type)
	entity.SetStatus(options.Status)

	return db.Model(entity).Updates(entity)
}

func (c *projectRepository) Delete(db *gorm.DB, dto *r.ProjectDto, entity []*e.ProjectEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *projectRepository) Reorder(db *gorm.DB, entity *e.ProjectEntity, position int) ([]*e.ProjectEntity, error) {
	var projects []*e.ProjectEntity
	db = db.Model(c.Model()).Where(`projects.organization_id = ?`, entity.OrganizationID)

	if entity.Order < position {
		db = db.Where(`projects.order > ?`, entity.Order).Where(`projects.order <= ?`, position)
	} else {
		db = db.Where(`projects.order < ?`, entity.Order).Where(`projects.order >= ?`, position)
	}

	res := db.Find(&projects)
	return projects, res.Error
}

var repository *projectRepository

func NewProjectRepository(db *gorm.DB) *ProjectRepositoryT {
	if repository == nil {
		repository = &projectRepository{}
	}

	return repositories.NewRepository(db, repository)
}
