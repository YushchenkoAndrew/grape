package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/context/dto/request"
	e "grape/src/context/entities"

	"gorm.io/gorm"
)

type ContextRelation string

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

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *contextRepository) sortBy(tx *gorm.DB, dto *r.ContextDto, _ []ContextRelation) {
}

func (c *contextRepository) Create(db *gorm.DB, dto *r.ContextDto, body interface{}, entity *e.ContextEntity) *gorm.DB {
	// db.Clauses(clause.Locking{Strength: clause.LockingStrengthShare, Table: clause.Table{Name: c.Model().TableName()}})
	// entity = body.(*e.ContextEntity)

	// if res := db.Create(entity); res.Error != nil {
	// 	return res
	// }

	// var order int64
	// db.Model(c.Model()).
	// 	Select(`MAX(attachments.order) AS "order"`).
	// 	Where(`attachments.attachable_id = ? AND attachments.attachable_type = ?`, entity.AttachableID, entity.AttachableType).
	// 	Scan(&order)

	// return db.Model(c.Model()).Where("attachments.id = ?", entity.ID).Update("order", int(order)+1)
	return db
}

func (c *contextRepository) Update(db *gorm.DB, dto *r.ContextDto, body interface{}, entity *e.ContextEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *contextRepository) Delete(db *gorm.DB, dto *r.ContextDto, entity []*e.ContextEntity) *gorm.DB {
	// for _, attachment := range entity {
	// 	var attachments []*e.ContextEntity
	// 	res := db.Model(c.Model()).
	// 		Where(`attachments.attachable_id = ? AND attachments.attachable_type = ?`, attachment.AttachableID, attachment.AttachableType).
	// 		Where(`attachments.order > ?`, attachment.Order).
	// 		Find(&attachments)

	// 	if res.Error != nil {
	// 		return res
	// 	}

	// 	if len(attachments) == 0 {
	// 		continue
	// 	}

	// 	lo.ForEach(attachments, func(e *e.ContextEntity, _ int) { e.Order -= 1 })
	// 	if res := db.Model(c.Model()).Save(attachments); res.Error != nil {
	// 		return res
	// 	}
	// }

	// return db.Model(c.Model()).Delete(entity)
	return db
}

func (c *contextRepository) Reorder(db *gorm.DB, entity *e.ContextEntity, position int) ([]*e.ContextEntity, error) {
	return []*e.ContextEntity{}, nil
}

var repository *contextRepository

func NewContextRepository(db *gorm.DB) *ContextRepositoryT {
	if repository == nil {
		repository = &contextRepository{}
	}

	return repositories.NewRepository(db, repository)
}
