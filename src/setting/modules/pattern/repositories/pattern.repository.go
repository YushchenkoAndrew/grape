package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/setting/modules/pattern/dto/request"
	e "grape/src/setting/modules/pattern/entities"
	"grape/src/setting/modules/pattern/types"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PatternRelation string

type PatternRepositoryT = repositories.CommonRepository[*e.PatternEntity, *r.PatternDto, PatternRelation]

type patternRepository struct {
}

func (c *patternRepository) Model() *e.PatternEntity {
	return e.NewPatternEntity()
}

func (c *patternRepository) Build(db *gorm.DB, dto *r.PatternDto, relations ...PatternRelation) *gorm.DB {
	tx := db.Model(c.Model()).Where(`patterns.organization_id = ?`, dto.CurrentUser.Organization.ID)

	required := c.applyFilter(tx, dto, []PatternRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *patternRepository) applyFilter(tx *gorm.DB, dto *r.PatternDto, relations []PatternRelation) []PatternRelation {
	if len(dto.PatternIds) != 0 {
		tx.Where(`patterns.uuid IN ?`, dto.PatternIds)
	}

	if len(dto.Colors) != 0 {
		tx.Where(`patterns.colors IN ?`, dto.Colors)
	}

	if len(dto.Modes) != 0 {
		tx.Where(`patterns.mode IN ?`, lo.Map(dto.Modes, func(str string, _ int) types.PatternColorModeEnum {
			return types.Fill.Value(str)
		}))
	}

	return relations
}

func (c *patternRepository) attachRelations(tx *gorm.DB, _ *r.PatternDto, relations []PatternRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *patternRepository) sortBy(tx *gorm.DB, dto *r.PatternDto, _ []PatternRelation) {
	switch dto.SortBy {
	case "name", "order", "created_at":
		tx.Order(repositories.NewSortBy(c.Model().TableName(), dto.SortBy, dto.Direction))

	default:
		tx.Order(repositories.NewSortBy(c.Model().TableName(), "id", dto.Direction))
	}

}

func (c *patternRepository) Create(db *gorm.DB, dto *r.PatternDto, body interface{}, entity *e.PatternEntity) *gorm.DB {
	options := body.(*r.PatternCreateDto)

	var data types.PatternOptionsType
	copier.Copy(&data, options.Options)
	entity.SetOptions(&data)

	entity.Organization = &dto.CurrentUser.Organization

	return db.Create(entity)
}

func (c *patternRepository) Update(db *gorm.DB, dto *r.PatternDto, body interface{}, entity *e.PatternEntity) *gorm.DB {
	options := body.(*r.PatternUpdateDto)

	if options.Options != nil {
		var data types.PatternOptionsType
		copier.CopyWithOption(&data, options.Options, copier.Option{IgnoreEmpty: true})
		entity.SetOptions(&data)
	}

	return db.Model(entity).Updates(entity)
}

func (c *patternRepository) Delete(db *gorm.DB, dto *r.PatternDto, entity *e.PatternEntity) *gorm.DB {
	// TODO: Add transaction with recursive delete related entities
	return db.Model(c.Model()).Delete(entity)
}

var repository *patternRepository

func NewPatternRepository(db *gorm.DB) *PatternRepositoryT {
	if repository == nil {
		repository = &patternRepository{}
	}

	return repositories.NewRepository(db, repository)
}
