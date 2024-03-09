package repositories

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
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
	Create(*Dto, *Entity) *gorm.DB
}

type CommonRepository[Entity CommonEntity, Dto CommonDtoT, Relations any] struct {
	handler CommonRepositoryT[Dto, Entity, Relations]
}

func (c *CommonRepository[Entity, Dto, Relations]) GetOne(dto *Dto, relations ...Relations) *Entity {
	var result []Entity
	c.handler.Build(dto, relations...).Limit(1).Find(&result)
	if len(result) == 0 {
		return nil
	}

	return &result[0]
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAll(dto *Dto, relations ...Relations) []Entity {
	var result []Entity
	c.handler.Build(dto, relations...).Find(&result)
	return result
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAllPage(dto Dto, relations ...Relations) (int, []Entity) {
	var cnt int64
	var result []Entity

	if dto.Limit() == 0 {
		return 0, result
	}

	c.handler.Build(&dto, relations...).Count(&cnt).Offset(dto.Offset()).Limit(dto.Limit()).Find(&result)
	return int(cnt), result
}

func (c *CommonRepository[Entity, Dto, Relations]) Create(dto *Dto, body interface{}) (*Entity, error) {
	entity := c.handler.Model()
	copier.Copy(&entity, body)

	if tx := c.handler.Create(dto, entity); tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil
}

func NewRepository[Entity CommonEntity, Dto CommonDtoT, Relations any](handler CommonRepositoryT[Dto, Entity, Relations]) *CommonRepository[Entity, Dto, Relations] {
	return &CommonRepository[Entity, Dto, Relations]{handler}
}
