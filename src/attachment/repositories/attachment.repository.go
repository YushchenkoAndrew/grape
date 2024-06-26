package repositories

import (
	r "grape/src/attachment/dto/request"
	e "grape/src/attachment/entities"
	"grape/src/common/repositories"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttachmentRelation string

type AttachmentRepositoryT = repositories.CommonRepository[*e.AttachmentEntity, *r.AttachmentDto, AttachmentRelation]

type attachmentRepository struct {
}

func (c *attachmentRepository) Model() *e.AttachmentEntity {
	return e.NewAttachmentEntity()
}

func (c *attachmentRepository) Build(db *gorm.DB, dto *r.AttachmentDto, relations ...AttachmentRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []AttachmentRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *attachmentRepository) applyFilter(tx *gorm.DB, dto *r.AttachmentDto, relations []AttachmentRelation) []AttachmentRelation {
	if len(dto.AttachmentIds) != 0 {
		tx.Where(`attachments.uuid IN ?`, dto.AttachmentIds)
	}

	return relations
}

func (c *attachmentRepository) attachRelations(tx *gorm.DB, _ *r.AttachmentDto, relations []AttachmentRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *attachmentRepository) sortBy(tx *gorm.DB, dto *r.AttachmentDto, _ []AttachmentRelation) {
}

func (c *attachmentRepository) Create(db *gorm.DB, dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	db.Clauses(clause.Locking{Strength: clause.LockingStrengthShare, Table: clause.Table{Name: c.Model().TableName()}})
	entity = body.(*e.AttachmentEntity)

	var order int64
	db.Model(c.Model()).
		Select(`COALESCE(MAX(attachments.order), 0) AS "order"`).
		Where(`attachments.attachable_id = ? AND attachments.attachable_type = ?`, entity.AttachableID, entity.AttachableType).
		Scan(&order)

	if res := db.Create(entity); res.Error != nil {
		return res
	}

	return db.Model(c.Model()).Where("attachments.id = ?", entity.ID).Update("order", int(order)+1)
}

func (c *attachmentRepository) Update(db *gorm.DB, dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *attachmentRepository) Delete(db *gorm.DB, dto *r.AttachmentDto, entity []*e.AttachmentEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *attachmentRepository) Reorder(db *gorm.DB, entity *e.AttachmentEntity, position int) ([]*e.AttachmentEntity, error) {
	var attachments []*e.AttachmentEntity
	db = db.Model(c.Model()).
		Where(`attachments.attachable_id = ? AND attachments.attachable_type = ?`, entity.AttachableID, entity.AttachableType)

	if entity.Order < position {
		db = db.Where(`attachments.order > ?`, entity.Order).Where(`attachments.order <= ?`, position)
	} else {
		db = db.Where(`attachments.order < ?`, entity.Order).Where(`attachments.order >= ?`, position)
	}

	res := db.Find(&attachments)
	return attachments, res.Error
}

var repository *attachmentRepository

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepositoryT {
	if repository == nil {
		repository = &attachmentRepository{}
	}

	return repositories.NewRepository(db, repository)
}
