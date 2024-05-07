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
	if len(dto.Statuses) != 0 {
		tx.Where(`stages.status IN ?`, lo.Map(dto.Statuses, func(str string, _ int) t.StatusEnum {
			return t.Active.Value(str)
		}))
	}

	return relations
}

func (c *stageRepository) attachRelations(tx *gorm.DB, _ *r.StageDto, relations []StageRelation) {
	for _, r := range relations {
		switch r {
		case Tasks:
			tx.Preload(string(r)).
				Preload(fmt.Sprintf("%s.Owner", r)).
				Preload(fmt.Sprintf("%s.Links", r), func(db *gorm.DB) *gorm.DB { return db.Order("Links.order ASC") }).
				Preload(fmt.Sprintf("%s.Attachments", r), func(db *gorm.DB) *gorm.DB { return db.Order("Attachments.order ASC") }).
				Preload(fmt.Sprintf("%s.Contexts", r), func(db *gorm.DB) *gorm.DB { return db.Order("Contexts.order ASC") }).
				Preload(fmt.Sprintf("%s.Contexts.ContextFields", r), func(db *gorm.DB) *gorm.DB { return db.Order("ContextFields.order ASC") })

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
	c.Build(db, dto).Select(`MAX(stages.order) AS "order"`).Scan(&order)

	entity.Order = int(order) + 1
	entity.Organization = &dto.CurrentUser.Organization

	return db.Create(entity)
}

func (c *stageRepository) Update(db *gorm.DB, dto *r.StageDto, body interface{}, entity *e.StageEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *stageRepository) Delete(db *gorm.DB, dto *r.StageDto, entity []*e.StageEntity) *gorm.DB {
	// // TODO: Add transaction with recursive delete related entities

	// TODO: Impl this
	// for _, project := range entity {
	// 	var projects []*e.StageEntity
	// 	res := db.Model(c.Model()).
	// 		Where(`projects.organization_id = ?`, project.OrganizationID).
	// 		Where(`projects.order > ?`, project.Order).
	// 		Find(&projects)

	// 	if res.Error != nil {
	// 		return res
	// 	}

	// 	if len(projects) == 0 {
	// 		continue
	// 	}

	// 	lo.ForEach(projects, func(e *entities.StageEntity, _ int) { e.Order -= 1 })
	// 	if res := db.Model(c.Model()).Save(projects); res.Error != nil {
	// 		return res
	// 	}
	// }

	return db.Model(c.Model()).Delete(entity)
}

func (c *stageRepository) Reorder(db *gorm.DB, entity *e.StageEntity, position int) ([]*e.StageEntity, error) {
	var stages []*e.StageEntity
	db = db.Model(c.Model()).Where(`projects.organization_id = ?`, entity.OrganizationID)

	if entity.Order < position {
		db = db.Where(`projects.order > ?`, entity.Order).Where(`projects.order <= ?`, position)
	} else {
		db = db.Where(`projects.order < ?`, entity.Order).Where(`projects.order >= ?`, position)
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
