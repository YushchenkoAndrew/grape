package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/tag/dto/request"
	e "grape/src/tag/entities"

	"gorm.io/gorm"
)

type TagRelation string

type TagRepositoryT = repositories.CommonRepository[*e.TagEntity, *r.TagDto, TagRelation]

type tagRepository struct {
}

func (c *tagRepository) Model() *e.TagEntity {
	return e.NewTagEntity()
}

func (c *tagRepository) Build(db *gorm.DB, dto *r.TagDto, relations ...TagRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []TagRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *tagRepository) applyFilter(tx *gorm.DB, dto *r.TagDto, relations []TagRelation) []TagRelation {
	if len(dto.TagIds) != 0 {
		tx.Where(`tags.uuid IN ?`, dto.TagIds)
	}

	return relations
}

func (c *tagRepository) attachRelations(tx *gorm.DB, _ *r.TagDto, relations []TagRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *tagRepository) sortBy(tx *gorm.DB, dto *r.TagDto, _ []TagRelation) {
}

func (c *tagRepository) Create(db *gorm.DB, dto *r.TagDto, body interface{}, entity *e.TagEntity) *gorm.DB {
	return db.Model(c.Model()).Create(body.(*e.TagEntity))
}

func (c *tagRepository) Update(db *gorm.DB, dto *r.TagDto, body interface{}, entity *e.TagEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *tagRepository) Delete(db *gorm.DB, dto *r.TagDto, entity []*e.TagEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *tagRepository) Reorder(db *gorm.DB, entity *e.TagEntity, position int) ([]*e.TagEntity, error) {
	return []*e.TagEntity{}, nil
}

var repository *tagRepository

func NewTagRepository(db *gorm.DB) *TagRepositoryT {
	if repository == nil {
		repository = &tagRepository{}
	}

	return repositories.NewRepository(db, repository)
}
