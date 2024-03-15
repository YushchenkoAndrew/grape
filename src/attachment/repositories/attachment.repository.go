package repositories

import (
	r "grape/src/attachment/dto/request"
	e "grape/src/attachment/entities"
	"grape/src/common/repositories"

	"gorm.io/gorm"
)

type AttachmentRelation string

type AttachmentRepositoryT = repositories.CommonRepository[*e.AttachmentEntity, *r.AttachmentDto, AttachmentRelation]

type attachmentRepository struct {
	db *gorm.DB
}

func (c *attachmentRepository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return c.db
}

func (c *attachmentRepository) Model() *e.AttachmentEntity {
	return e.NewAttachmentEntity()
}

func (c *attachmentRepository) Transaction(fc func(*gorm.DB) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		return fc(tx.Model(c.Model()))
	})
}

func (c *attachmentRepository) Build(db *gorm.DB, dto *r.AttachmentDto, relations ...AttachmentRelation) *gorm.DB {
	tx := c.conn(db).Model(c.Model())

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

func (c *attachmentRepository) Create(tx *gorm.DB, dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	return nil
}

func (c *attachmentRepository) Update(tx *gorm.DB, dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	return c.conn(tx).Model(entity).Updates(entity)
}

func (c *attachmentRepository) Delete(tx *gorm.DB, dto *r.AttachmentDto, entity *e.AttachmentEntity) *gorm.DB {
	return c.conn(tx).Model(c.Model()).Delete(entity)
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepositoryT {
	return repositories.NewRepository(&attachmentRepository{db})
}
