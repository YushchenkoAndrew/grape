package repositories

import (
	"grape/src/common/repositories"
	t "grape/src/common/types"
	r "grape/src/stage/dto/request"
	e "grape/src/stage/entities"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type TaskRelation string

// const (
// 	Organization StageRelation = "Organization"
// )

type TaskRepositoryT = repositories.CommonRepository[*e.TaskEntity, *r.TaskDto, TaskRelation]

type taskRepository struct {
}

func (c *taskRepository) Model() *e.TaskEntity {
	return e.NewTaskEntity()
}

func (c *taskRepository) Build(db *gorm.DB, dto *r.TaskDto, relations ...TaskRelation) *gorm.DB {
	tx := db.Model(c.Model()).Where(`tasks.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []TaskRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *taskRepository) applyFilter(tx *gorm.DB, dto *r.TaskDto, relations []TaskRelation) []TaskRelation {
	if len(dto.Statuses) != 0 {
		tx.Where(`tasks.status IN ?`, lo.Map(dto.Statuses, func(str string, _ int) t.StatusEnum {
			return t.Active.Value(str)
		}))
	}

	if len(dto.StageIds) != 0 {
		tx.Where(`tasks.stage_id IN (SELECT id FROM stages WHERE uuid IN ?)`, dto.StageIds)
	}

	return relations
}

func (c *taskRepository) attachRelations(tx *gorm.DB, _ *r.TaskDto, relations []TaskRelation) {
	for _, r := range relations {
		switch r {
		default:
			tx.Joins(string(r))
		}
	}
}

func (c *taskRepository) sortBy(tx *gorm.DB, dto *r.TaskDto, _ []TaskRelation) {
	switch dto.SortBy {
	case "":
		return

	case "name", "order", "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *taskRepository) Create(db *gorm.DB, dto *r.TaskDto, body interface{}, entity *e.TaskEntity) *gorm.DB {
	var order int64
	dto.SortBy = ""
	c.Build(db, dto).Select(`MAX(tasks.order) AS "order"`).Scan(&order)
	options := body.(*r.TaskCreateDto)

	entity.Order = int(order) + 1
	entity.StageID = options.Stage.ID
	entity.Organization = &dto.CurrentUser.Organization

	return db.Create(entity)
}

func (c *taskRepository) Update(db *gorm.DB, dto *r.TaskDto, body interface{}, entity *e.TaskEntity) *gorm.DB {
	// options := body.(*r.ProjectUpdateDto)

	// // entity.SetType(options.Type)
	// entity.SetStatus(options.Status)

	// return db.Model(entity).Updates(entity)

	return db
}

func (c *taskRepository) Delete(db *gorm.DB, dto *r.TaskDto, entity []*e.TaskEntity) *gorm.DB {
	// // TODO: Add transaction with recursive delete related entities

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

	// return db.Model(c.Model()).Delete(entity)
	return db
}

func (c *taskRepository) Reorder(db *gorm.DB, entity *e.TaskEntity, position int) ([]*e.TaskEntity, error) {
	var tasks []*e.TaskEntity
	// TODO: Think about how to move it to the next stage !!!
	// db = db.Model(c.Model()).Where(`tasks.organization_id = ?`, entity.OrganizationID)

	// if entity.Order < position {
	// 	db = db.Where(`projects.order > ?`, entity.Order).Where(`projects.order <= ?`, position)
	// } else {
	// 	db = db.Where(`projects.order < ?`, entity.Order).Where(`projects.order >= ?`, position)
	// }

	// res := db.Find(&projects)
	return tasks, nil
}

var task_repository *taskRepository

func NewTaskRepository(db *gorm.DB) *TaskRepositoryT {
	if repository == nil {
		task_repository = &taskRepository{}
	}

	return repositories.NewRepository(db, task_repository)
}