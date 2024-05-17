package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/statistic/dto/request"
	e "grape/src/statistic/entities"

	"gorm.io/gorm"
)

type StatisticRelation string

const (
	Project StatisticRelation = "Project"
)

type StatisticRepositoryT = repositories.CommonRepository[*e.StatisticEntity, *r.StatisticDto, StatisticRelation]

type statisticRepository struct {
}

func (c *statisticRepository) Model() *e.StatisticEntity {
	return e.NewStatisticEntity()
}

func (c *statisticRepository) Build(db *gorm.DB, dto *r.StatisticDto, relations ...StatisticRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []StatisticRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *statisticRepository) applyFilter(tx *gorm.DB, dto *r.StatisticDto, relations []StatisticRelation) []StatisticRelation {
	if len(dto.ProjectIds) != 0 {
		tx.Where(`projects.uuid IN ?`, dto.ProjectIds)
	}

	return relations
}

func (c *statisticRepository) attachRelations(tx *gorm.DB, _ *r.StatisticDto, relations []StatisticRelation) {
	for _, r := range relations {
		switch r {
		default:
			tx.Joins(string(r))
		}
	}
}

func (c *statisticRepository) sortBy(tx *gorm.DB, dto *r.StatisticDto, _ []StatisticRelation) {}

func (c *statisticRepository) Create(db *gorm.DB, dto *r.StatisticDto, body interface{}, entity *e.StatisticEntity) *gorm.DB {
	return nil
}

func (c *statisticRepository) Update(db *gorm.DB, dto *r.StatisticDto, body interface{}, entity *e.StatisticEntity) *gorm.DB {
	options := body.(*r.StatisticUpdateDto)

	switch options.Kind {
	case "view":
		entity.Views += 1

	case "click":
		entity.Clicks += 1

	case "media":
		entity.Media += 1
	}

	return db.Model(entity).Updates(entity)
}

func (c *statisticRepository) Delete(db *gorm.DB, dto *r.StatisticDto, entity []*e.StatisticEntity) *gorm.DB {
	return nil
}

func (c *statisticRepository) Reorder(db *gorm.DB, entity *e.StatisticEntity, position int) ([]*e.StatisticEntity, error) {
	return []*e.StatisticEntity{}, nil
}

var repository *statisticRepository

func NewStatisticRepository(db *gorm.DB) *StatisticRepositoryT {
	if repository == nil {
		repository = &statisticRepository{}
	}

	return repositories.NewRepository(db, repository)
}
