package repositories

import (
	r "grape/src/attachment/dto/request"
	e "grape/src/attachment/entities"
	"grape/src/common/repositories"

	"gorm.io/gorm"
)

type AttachmentRelation string

type AttachmentRepositoryT = repositories.CommonRepository[*e.AttachmentEntity, *r.AttachmentDto, AttachmentRelation]

type locationRepository struct {
	db *gorm.DB
}

func (c *locationRepository) Model() *e.AttachmentEntity {
	return e.NewAttachmentEntity()
}

func (c *locationRepository) Transaction(_ func(*gorm.DB) error) error {
	return nil
}

func (c *locationRepository) Build(dto *r.AttachmentDto, relations ...AttachmentRelation) *gorm.DB {
	tx := c.db.Model(c.Model())

	required := c.applyFilter(tx, dto, []AttachmentRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *locationRepository) applyFilter(tx *gorm.DB, dto *r.AttachmentDto, relations []AttachmentRelation) []AttachmentRelation {
	if len(dto.AttachmentIds) != 0 {
		tx.Where(`attachments.uuid IN ?`, dto.AttachmentIds)
	}

	return relations
}

func (c *locationRepository) attachRelations(tx *gorm.DB, _ *r.AttachmentDto, relations []AttachmentRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *locationRepository) sortBy(tx *gorm.DB, dto *r.AttachmentDto, _ []AttachmentRelation) {
}

func (c *locationRepository) Create(dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	return nil
}

func (c *locationRepository) Update(dto *r.AttachmentDto, body interface{}, entity *e.AttachmentEntity) *gorm.DB {
	return nil
}

func (c *locationRepository) Delete(dto *r.AttachmentDto, entity *e.AttachmentEntity) *gorm.DB {
	return nil
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepositoryT {
	return repositories.NewRepository(&locationRepository{db})
}
