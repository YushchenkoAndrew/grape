package repositories

import (
	"grape/src/common/repositories"
	r "grape/src/customer/dto/request"
	e "grape/src/customer/entities"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type LocationRelation string

const (
	Network LocationRelation = "Network"
)

type LocationRepositoryT = repositories.CommonRepository[*e.LocationEntity, *r.LocationDto, LocationRelation]

type locationRepository struct {
}

func (c *locationRepository) Model() *e.LocationEntity {
	return e.NewLocationEntity()
}

func (c *locationRepository) Build(db *gorm.DB, dto *r.LocationDto, relations ...LocationRelation) *gorm.DB {
	tx := db.Model(c.Model())

	required := c.applyFilter(tx, dto, []LocationRelation{})
	c.attachRelations(tx, dto, append(relations, required...))
	c.sortBy(tx, dto, append(relations, required...))

	return tx
}

func (c *locationRepository) applyFilter(tx *gorm.DB, dto *r.LocationDto, relations []LocationRelation) []LocationRelation {
	if len(dto.IP) != 0 {
		relations = append(relations, Network)
		tx.Where(lo.Reduce(dto.IP, func(acc *gorm.DB, curr string, _ int) *gorm.DB {
			return acc.Where(`network >>= ?::inet`, curr)
		}, tx))
	}

	return relations
}

func (c *locationRepository) attachRelations(tx *gorm.DB, _ *r.LocationDto, relations []LocationRelation) {
	for _, r := range relations {
		switch r {

		default:
			tx.Joins(string(r))
		}
	}
}

func (c *locationRepository) sortBy(tx *gorm.DB, dto *r.LocationDto, _ []LocationRelation) {
}

func (c *locationRepository) Create(tx *gorm.DB, dto *r.LocationDto, body interface{}, entity *e.LocationEntity) *gorm.DB {
	return nil
}

func (c *locationRepository) Update(tx *gorm.DB, dto *r.LocationDto, body interface{}, entity *e.LocationEntity) *gorm.DB {
	return nil
}

func (c *locationRepository) Delete(tx *gorm.DB, dto *r.LocationDto, entity []*e.LocationEntity) *gorm.DB {
	return nil
}

func (c *locationRepository) Reorder(db *gorm.DB, entity *e.LocationEntity, position int) ([]*e.LocationEntity, error) {
	return []*e.LocationEntity{}, nil
}

var repository *locationRepository

func NewLocationRepository(db *gorm.DB) *LocationRepositoryT {
	if repository == nil {
		repository = &locationRepository{}
	}

	return repositories.NewRepository(db, repository)
}
