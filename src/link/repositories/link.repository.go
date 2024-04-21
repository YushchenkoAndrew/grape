package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/link/dto/request"
	e "grape/src/link/entities"

	"gorm.io/gorm"
)

type LinkRelation string

type LinkRepositoryT = repositories.CommonRepository[*e.LinkEntity, *r.LinkDto, LinkRelation]

type linkRepository struct {
}

func (c *linkRepository) Model() *e.LinkEntity {
	return e.NewLinkEntity()
}

func (c *linkRepository) Build(db *gorm.DB, dto *r.LinkDto, relations ...LinkRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []LinkRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *linkRepository) applyFilter(tx *gorm.DB, dto *r.LinkDto, relations []LinkRelation) []LinkRelation {
	if len(dto.LinkIds) != 0 {
		tx.Where(`links.uuid IN ?`, dto.LinkIds)
	}

	return relations
}

func (c *linkRepository) attachRelations(tx *gorm.DB, _ *r.LinkDto, relations []LinkRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *linkRepository) sortBy(tx *gorm.DB, dto *r.LinkDto, _ []LinkRelation) {
}

func (c *linkRepository) Create(db *gorm.DB, dto *r.LinkDto, body interface{}, entity *e.LinkEntity) *gorm.DB {
	return nil
}

func (c *linkRepository) Update(db *gorm.DB, dto *r.LinkDto, body interface{}, entity *e.LinkEntity) *gorm.DB {
	return db.Model(entity).Updates(entity)
}

func (c *linkRepository) Delete(db *gorm.DB, dto *r.LinkDto, entity []*e.LinkEntity) *gorm.DB {
	return db.Model(c.Model()).Delete(entity)
}

func (c *linkRepository) Reorder(db *gorm.DB, entity *e.LinkEntity, position int) error {
	return nil
}

var repository *linkRepository

func NewLinkRepository(db *gorm.DB) *LinkRepositoryT {
	if repository == nil {
		repository = &linkRepository{}
	}

	return repositories.NewRepository(db, repository)
}
