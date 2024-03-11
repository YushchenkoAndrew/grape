package repositories

import (
	"grape/src/common/types"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonDtoT interface {
	Offset() int
	Limit() int
}

type CommonEntity interface {
	Create()
	Update()
}

type CommonRepositoryT[Dto CommonDtoT, Entity CommonEntity, Relations any] interface {
	Model() *Entity
	Build(*Dto, ...Relations) *gorm.DB
	Create(*Dto, interface{}, *Entity) *gorm.DB
	Update(*Dto, interface{}, *Entity) *gorm.DB
	Delete(*Dto, *Entity) *gorm.DB
}

type CommonRepository[Entity CommonEntity, Dto CommonDtoT, Relations any] struct {
	handler CommonRepositoryT[Dto, Entity, Relations]
}

func (c *CommonRepository[Entity, Dto, Relations]) GetOne(dto *Dto, relations ...Relations) (*Entity, error) {
	var result []Entity
	if tx := c.handler.Build(dto, relations...).Limit(1).Find(&result); tx.Error != nil {
		return nil, tx.Error
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAll(dto *Dto, relations ...Relations) ([]Entity, error) {
	var result []Entity
	if tx := c.handler.Build(dto, relations...).Find(&result); tx.Error != nil {
		return nil, tx.Error
	}

	return result, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAllPage(dto *Dto, relations ...Relations) (int, []Entity, error) {
	var cnt int64
	var result []Entity

	if (*dto).Limit() == 0 {
		return 0, result, nil
	}

	if tx := c.handler.Build(dto, relations...).Count(&cnt); tx.Error != nil {
		return 0, nil, tx.Error
	}

	if tx := c.handler.Build(dto, relations...).Offset((*dto).Offset()).Limit((*dto).Limit()).Find(&result); tx.Error != nil {
		return 0, nil, tx.Error
	}

	return int(cnt), result, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Create(dto *Dto, body interface{}) (*Entity, error) {
	entity := c.handler.Model()
	copier.Copy(&entity, body)
	(*entity).Create()

	if tx := c.handler.Create(dto, body, entity); tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Update(dto *Dto, body interface{}, entity *Entity) (*Entity, error) {
	copier.CopyWithOption(&entity, body, copier.Option{IgnoreEmpty: true})
	(*entity).Update()

	if tx := c.handler.Update(dto, body, entity); tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Delete(dto *Dto, entity *Entity) error {
	if tx := c.handler.Delete(dto, entity); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func NewSortBy(alias, column string, direction string) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: column, Table: alias}, Desc: types.Asc.Value(direction).Bool()}
}

func NewRepository[Entity CommonEntity, Dto CommonDtoT, Relations any](handler CommonRepositoryT[Dto, Entity, Relations]) *CommonRepository[Entity, Dto, Relations] {
	return &CommonRepository[Entity, Dto, Relations]{handler}
}
