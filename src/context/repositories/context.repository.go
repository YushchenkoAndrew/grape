package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/context/dto/request"
	e "grape/src/context/entities"

	"gorm.io/gorm"
)

type ContextRelation string

const (
	ContextFields ContextRelation = "ContextFields"
)

type ContextRepositoryT = repositories.CommonRepository[*e.ContextEntity, *r.ContextDto, ContextRelation]

type contextRepository struct {
}

func (c *contextRepository) Model() *e.ContextEntity {
	return e.NewContextEntity()
}

func (c *contextRepository) Build(db *gorm.DB, dto *r.ContextDto, relations ...ContextRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []ContextRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *contextRepository) applyFilter(tx *gorm.DB, dto *r.ContextDto, relations []ContextRelation) []ContextRelation {
	if len(dto.ContextIds) != 0 {
		tx.Where(`contexts.uuid IN ?`, dto.ContextIds)
	}

	return relations
}

func (c *contextRepository) attachRelations(tx *gorm.DB, _ *r.ContextDto, relations []ContextRelation) {
	for _, r := range relations {
		switch r {
		case ContextFields:
			tx.Preload(string(r))

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *contextRepository) sortBy(tx *gorm.DB, dto *r.ContextDto, _ []ContextRelation) {
}

func (c *contextRepository) Create(db *gorm.DB, dto *r.ContextDto, body interface{}, entity *e.ContextEntity) *gorm.DB {
	entity = body.(*e.ContextEntity)

	var order int64
	db.Model(c.Model()).
		Select(`COALESCE(MAX(contexts.order), 0) AS "order"`).
		Where(`contexts.contextable_id = ? AND contexts.contextable_type = ?`, entity.ContextableID, entity.ContextableType).
		Scan(&order)

	if res := db.Create(entity); res.Error != nil {
		return res
	}

	return db.Model(c.Model()).Where("contexts.id = ?", entity.ID).Update("order", int(order)+1)
}

func (c *contextRepository) Update(db *gorm.DB, dto *r.ContextDto, body interface{}, entity *e.ContextEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *contextRepository) Delete(db *gorm.DB, dto *r.ContextDto, entity []*e.ContextEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *contextRepository) Reorder(db *gorm.DB, entity *e.ContextEntity, position int) ([]*e.ContextEntity, error) {
	var contexts []*e.ContextEntity
	db = db.Model(c.Model()).
		Where(`contexts.contextable_id = ? AND contexts.contextable_type = ?`, entity.ContextableID, entity.ContextableType)

	if entity.Order < position {
		db = db.Where(`contexts.order > ?`, entity.Order).Where(`contexts.order <= ?`, position)
	} else {
		db = db.Where(`contexts.order < ?`, entity.Order).Where(`contexts.order >= ?`, position)
	}

	res := db.Find(&contexts)
	return contexts, res.Error
}

var repository *contextRepository

func NewContextRepository(db *gorm.DB) *ContextRepositoryT {
	if repository == nil {
		repository = &contextRepository{}
	}

	return repositories.NewRepository(db, repository)
}
