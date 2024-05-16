package repositories

import (
	"fmt"
	"grape/src/common/repositories"
	t "grape/src/common/types"
	r "grape/src/stage/dto/request"
	e "grape/src/stage/entities"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StageRelation string

const (
	Organization StageRelation = "Organization"
	Tasks        StageRelation = "Tasks"
)

type StageRepositoryT = repositories.CommonRepository[*e.StageEntity, *r.StageDto, StageRelation]

type stageRepository struct {
}

func (c *stageRepository) Model() *e.StageEntity {
	return e.NewStageEntity()
}

func (c *stageRepository) Build(db *gorm.DB, dto *r.StageDto, relations ...StageRelation) *gorm.DB {
	tx := db.Model(c.Model()).Where(`stages.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []StageRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *stageRepository) applyFilter(tx *gorm.DB, dto *r.StageDto, relations []StageRelation) []StageRelation {
	if len(dto.StageIds) != 0 {
		tx.Where(`stages.uuid IN ?`, dto.StageIds)
	}

	if len(dto.Statuses) != 0 {
		tx.Where(`stages.status IN ?`, lo.Map(dto.Statuses, func(str string, _ int) t.StatusEnum {
			return t.Active.Value(str)
		}))
	}

	return relations
}

func (c *stageRepository) attachRelations(tx *gorm.DB, dto *r.StageDto, relations []StageRelation) {
	for _, r := range relations {
		switch r {
		case Tasks:
			tx.Preload(string(r), func(db *gorm.DB) *gorm.DB {
				if len(dto.Query) != 0 {
					query := "%" + dto.Query + "%"
					db.Where(
						db.Session(&gorm.Session{NewDB: true}).Where(`tasks.name ILIKE ?`, query).
							Or(`(SELECT COUNT(t.id) FROM tags t WHERE t.taggable_id = tasks.id AND t.taggable_type = ? AND t.name ILIKE ?) > 0`, string(r), query),
					)
				}

				return db.Order("tasks.order ASC")
			}).
				Preload(fmt.Sprintf("%s.Owner", r)).
				Preload(fmt.Sprintf("%s.Tags", r)).
				Preload(fmt.Sprintf("%s.Links", r), func(db *gorm.DB) *gorm.DB { return db.Order("links.order ASC") }).
				Preload(fmt.Sprintf("%s.Attachments", r), func(db *gorm.DB) *gorm.DB { return db.Order("attachments.order ASC") }).
				Preload(fmt.Sprintf("%s.Contexts", r), func(db *gorm.DB) *gorm.DB { return db.Order("contexts.order ASC") }).
				Preload(fmt.Sprintf("%s.Contexts.ContextFields", r), func(db *gorm.DB) *gorm.DB { return db.Order("context_fields.order ASC") })

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *stageRepository) sortBy(tx *gorm.DB, dto *r.StageDto, _ []StageRelation) {
	switch dto.SortBy {
	case "":
		return

	case "name", "order", "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *stageRepository) Create(db *gorm.DB, dto *r.StageDto, body interface{}, entity *e.StageEntity) *gorm.DB {
	var order int64
	dto.SortBy = ""
	c.Build(db, dto).Select(`COALESCE(MAX(stages.order), 0) AS "order"`).Scan(&order)

	entity.Order = int(order) + 1
	entity.Organization = &dto.CurrentUser.Organization

	return db.Create(entity)
}

func (c *stageRepository) Update(db *gorm.DB, dto *r.StageDto, body interface{}, entity *e.StageEntity) *gorm.DB {
	options := body.(*r.StageUpdateDto)
	entity.SetStatus(options.Status)

	return db.Model(entity).Updates(entity)
}

func (c *stageRepository) Delete(db *gorm.DB, dto *r.StageDto, entity []*e.StageEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *stageRepository) Reorder(db *gorm.DB, entity *e.StageEntity, position int) ([]*e.StageEntity, error) {
	var stages []*e.StageEntity
	db = db.Model(c.Model()).Where(`stages.organization_id = ?`, entity.OrganizationID)

	if entity.Order < position {
		db = db.Where(`stages.order > ?`, entity.Order).Where(`stages.order <= ?`, position)
	} else {
		db = db.Where(`stages.order < ?`, entity.Order).Where(`stages.order >= ?`, position)
	}

	res := db.Find(&stages)
	return stages, res.Error
}

var repository *stageRepository

func NewStageRepository(db *gorm.DB) *StageRepositoryT {
	if repository == nil {
		repository = &stageRepository{}
	}

	return repositories.NewRepository(db, repository)
}
