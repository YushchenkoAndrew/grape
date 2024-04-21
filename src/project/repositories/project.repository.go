package repositories

import (
	"fmt"
	"grape/src/common/repositories"
	r "grape/src/project/dto/request"
	"grape/src/project/entities"
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
	Owner        ProjectRelation = "Owner"
	Palette      ProjectRelation = "Palette"
	Pattern      ProjectRelation = "Pattern"
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
		case Thumbnail:
			tx.Preload(string(r), "name ILIKE ?", "thumbnail%")

		case Redirect:
			tx.Preload(string(r), "name ILIKE ?", "redirect")

		case Attachments:
			tx.Preload(string(r), func(db *gorm.DB) *gorm.DB { return db.Order(fmt.Sprintf("%s.order ASC", string(r))) })

		case Links:
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
	c.Build(db, dto).Select(`MAX(projects.order) AS "order"`).Scan(&order)

	db.First(&entity.Pattern)
	db.First(&entity.Palette)

	entity.Order = int(order) + 1
	entity.Owner = dto.CurrentUser
	entity.Organization = &dto.CurrentUser.Organization
	entity.Statistic = st.NewStatisticEntity()

	entity.SetType(body.(*r.ProjectCreateDto).Type)
	return db.Create(entity)
}

func (c *projectRepository) Update(db *gorm.DB, dto *r.ProjectDto, body interface{}, entity *e.ProjectEntity) *gorm.DB {
	options := body.(*r.ProjectUpdateDto)

	if options.PaletteID != "" {
		db.First(&entity.Palette, "uuid = ?", options.PaletteID)
	}

	if options.PatternID != "" {
		db.First(&entity.Pattern, "uuid = ?", options.PatternID)
	}

	// entity.SetType(options.Type)
	entity.SetStatus(options.Status)

	return db.Model(entity).Updates(entity)
}

func (c *projectRepository) Delete(db *gorm.DB, dto *r.ProjectDto, entity []*e.ProjectEntity) *gorm.DB {
	// TODO: Add transaction with recursive delete related entities

	for _, project := range entity {
		var projects []*e.ProjectEntity
		res := db.Model(c.Model()).
			Where(`projects.organization_id = ?`, project.OrganizationID).
			Where(`projects.order > ?`, project.Order).
			Find(&projects)

		if res.Error != nil {
			return res
		}

		if len(projects) == 0 {
			continue
		}

		lo.ForEach(projects, func(e *entities.ProjectEntity, _ int) { e.Order -= 1 })
		if res := db.Model(c.Model()).Save(projects); res.Error != nil {
			return res
		}
	}

	return db.Model(c.Model()).Delete(entity)
}

func (c *projectRepository) Reorder(db *gorm.DB, entity *e.ProjectEntity, position int) error {
	var projects []*e.ProjectEntity
	db = db.Model(c.Model()).Where(`projects.organization_id = ?`, entity.OrganizationID)

	if entity.Order < position {
		db = db.Where(`projects.order > ?`, entity.Order).Where(`projects.order <= ?`, position)
	} else {
		db = db.Where(`projects.order < ?`, entity.Order).Where(`projects.order >= ?`, position)
	}

	if res := db.Find(&projects); res.Error != nil || len(projects) == 0 {
		return res.Error
	}

	for _, e := range projects {
		if entity.Order < position {
			e.Order -= 1
		} else {
			e.Order += 1
		}
	}

	entity.Order = position
	projects = append(projects, entity)

	res := db.Model(c.Model()).Save(projects)
	return res.Error
}

var repository *projectRepository

func NewProjectRepository(db *gorm.DB) *ProjectRepositoryT {
	if repository == nil {
		repository = &projectRepository{}
	}

	return repositories.NewRepository(db, repository)
}
