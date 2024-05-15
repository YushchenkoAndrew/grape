package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/context/dto/request"
	e "grape/src/context/entities"

	"gorm.io/gorm"
)

type ContextFieldRelation string

type ContextFieldRepositoryT = repositories.CommonRepository[*e.ContextFieldEntity, *r.ContextFieldDto, ContextFieldRelation]

type contextFieldRepository struct {
}

func (c *contextFieldRepository) Model() *e.ContextFieldEntity {
	return e.NewContextFieldEntity()
}

func (c *contextFieldRepository) Build(db *gorm.DB, dto *r.ContextFieldDto, relations ...ContextFieldRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []ContextFieldRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *contextFieldRepository) applyFilter(tx *gorm.DB, dto *r.ContextFieldDto, relations []ContextFieldRelation) []ContextFieldRelation {
	if len(dto.ContextFieldIds) != 0 {
		tx.Where(`context_fields.uuid IN ?`, dto.ContextFieldIds)
	}

	if len(dto.ContextIds) != 0 {
		tx.Where(`context_fields.context_id IN (SELECT id FROM contexts WHERE uuid IN ?)`, dto.ContextIds)
	}

	return relations
}

func (c *contextFieldRepository) attachRelations(tx *gorm.DB, _ *r.ContextFieldDto, relations []ContextFieldRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *contextFieldRepository) sortBy(tx *gorm.DB, dto *r.ContextFieldDto, _ []ContextFieldRelation) {
}

func (c *contextFieldRepository) Create(db *gorm.DB, dto *r.ContextFieldDto, body interface{}, entity *e.ContextFieldEntity) *gorm.DB {
	var order int64
	c.Build(db, dto).Select(`COALESCE(MAX(tasks.order), 0) AS "order"`).Scan(&order)
	options := body.(*r.ContextFieldCreateDto)

	entity.Order = int(order) + 1
	entity.ContextID = options.Context.ID
	entity.SetOptions(options.Options)

	return db.Create(entity)
}

func (c *contextFieldRepository) Update(db *gorm.DB, dto *r.ContextFieldDto, body interface{}, entity *e.ContextFieldEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *contextFieldRepository) Delete(db *gorm.DB, dto *r.ContextFieldDto, entity []*e.ContextFieldEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *contextFieldRepository) Reorder(db *gorm.DB, entity *e.ContextFieldEntity, position int) ([]*e.ContextFieldEntity, error) {
	var contexts []*e.ContextFieldEntity
	db = db.Model(c.Model()).
		Where(`context_fields.context_id = ?`, entity.ContextID)

	if entity.Order < position {
		db = db.Where(`context_fields.order > ?`, entity.Order).Where(`context_fields.order <= ?`, position)
	} else {
		db = db.Where(`context_fields.order < ?`, entity.Order).Where(`context_fields.order >= ?`, position)
	}

	res := db.Find(&contexts)
	return contexts, res.Error
}

var context_field *contextFieldRepository

func NewContextFieldRepository(db *gorm.DB) *ContextFieldRepositoryT {
	if context_field == nil {
		context_field = &contextFieldRepository{}
	}

	return repositories.NewRepository(db, context_field)
}
