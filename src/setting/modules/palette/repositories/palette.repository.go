package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/setting/modules/palette/dto/request"
	e "grape/src/setting/modules/palette/entities"

	"gorm.io/gorm"
)

type PaletteRelation string

type PaletteRepositoryT = repositories.CommonRepository[*e.PaletteEntity, *r.PaletteDto, PaletteRelation]

type paletteRepository struct {
}

func (c *paletteRepository) Model() *e.PaletteEntity {
	return e.NewPaletteEntity()
}

func (c *paletteRepository) Build(db *gorm.DB, dto *r.PaletteDto, relations ...PaletteRelation) *gorm.DB {
	tx := db.Model(c.Model()).Where(`palettes.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []PaletteRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *paletteRepository) applyFilter(tx *gorm.DB, dto *r.PaletteDto, relations []PaletteRelation) []PaletteRelation {
	if len(dto.PaletteIds) != 0 {
		tx.Where(`palettes.uuid IN ?`, dto.PaletteIds)
	}

	return relations
}

func (c *paletteRepository) attachRelations(tx *gorm.DB, _ *r.PaletteDto, relations []PaletteRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *paletteRepository) sortBy(tx *gorm.DB, dto *r.PaletteDto, _ []PaletteRelation) {
	switch dto.SortBy {
	case "name", "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *paletteRepository) Create(db *gorm.DB, dto *r.PaletteDto, body interface{}, entity *e.PaletteEntity) *gorm.DB {
	entity.Organization = &dto.CurrentUser.Organization
	return db.Create(entity)
}

func (c *paletteRepository) Update(db *gorm.DB, dto *r.PaletteDto, body interface{}, entity *e.PaletteEntity) *gorm.DB {
	// options := body.(*r.PaletteCreateDto)

	return db.Model(entity).Updates(entity)
}

func (c *paletteRepository) Delete(db *gorm.DB, dto *r.PaletteDto, entity *e.PaletteEntity) *gorm.DB {
	// TODO: Add transaction with recursive delete related entities
	return db.Model(c.Model()).Delete(entity)
}

var repository *paletteRepository

func NewPaletteRepository(db *gorm.DB) *PaletteRepositoryT {
	if repository == nil {
		repository = &paletteRepository{}
	}

	return repositories.NewRepository(db, repository)
}
