package repositories

import (
	"fmt"
	"grape/src/common/types"
	"reflect"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonDtoT interface {
	Offset() int
	Limit() int
	UUID() string
}

type CommonEntity interface {
	Create()
	Update()
	TableName() string
}

type CommonRepositoryT[Dto CommonDtoT, Entity CommonEntity, Relations any] interface {
	Model() Entity
	Transaction(func(*gorm.DB) error) error

	Build(*gorm.DB, Dto, ...Relations) *gorm.DB
	Create(*gorm.DB, Dto, interface{}, Entity) *gorm.DB
	Update(*gorm.DB, Dto, interface{}, Entity) *gorm.DB
	Delete(*gorm.DB, Dto, Entity) *gorm.DB
}

type CommonRepository[Entity CommonEntity, Dto CommonDtoT, Relations any] struct {
	handler CommonRepositoryT[Dto, Entity, Relations]
}

func (c *CommonRepository[Entity, Dto, Relations]) TableName() string {
	return c.handler.Model().TableName()
}

func (c *CommonRepository[Entity, Dto, Relations]) GetOne(dto Dto, relations ...Relations) (Entity, error) {
	var result []Entity
	if tx := c.handler.Build(nil, dto, relations...).Limit(1).Find(&result); tx.Error != nil {
		var e Entity
		return e, tx.Error
	}

	if len(result) == 0 {
		var e Entity
		return e, nil
	}

	return result[0], nil
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAll(dto Dto, relations ...Relations) ([]Entity, error) {
	var result []Entity
	if tx := c.handler.Build(nil, dto, relations...).Find(&result); tx.Error != nil {
		return nil, tx.Error
	}

	return result, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAllPage(dto Dto, relations ...Relations) (int, []Entity, error) {
	var cnt int64
	var result []Entity

	if dto.Limit() == 0 {
		return 0, result, nil
	}

	if tx := c.handler.Build(nil, dto, relations...).Count(&cnt); tx.Error != nil {
		return 0, nil, tx.Error
	}

	if tx := c.handler.Build(nil, dto, relations...).Offset(dto.Offset()).Limit(dto.Limit()).Find(&result); tx.Error != nil {
		return 0, nil, tx.Error
	}

	return int(cnt), result, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) ValidateEntityExistence(dto Dto, relations ...Relations) (Entity, error) {
	if uuid.Validate(dto.UUID()) != nil {
		var e Entity
		return e, fmt.Errorf("%s id is invalid", c.TableName())
	}

	if result, _ := c.GetOne(dto, relations...); !reflect.ValueOf(result).IsNil() {
		return result, nil
	}

	var e Entity
	return e, fmt.Errorf("%s not found", c.TableName())
}

func (c *CommonRepository[Entity, Dto, Relations]) Create(tx *gorm.DB, dto Dto, body interface{}) (*Entity, error) {
	entity := c.handler.Model()
	copier.Copy(&entity, body)
	entity.Create()

	if tx := c.handler.Create(tx, dto, body, entity); tx.Error != nil {
		return nil, tx.Error
	}

	return &entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Update(tx *gorm.DB, dto Dto, body interface{}, entity Entity) (Entity, error) {
	copier.CopyWithOption(&entity, body, copier.Option{IgnoreEmpty: true})
	entity.Update()

	if tx := c.handler.Update(tx, dto, body, entity); tx.Error != nil {
		var e Entity
		return e, tx.Error
	}

	return entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Delete(tx *gorm.DB, dto Dto, entity Entity) error {
	if tx := c.handler.Delete(tx, dto, entity); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Transaction(fc func(tx *gorm.DB) error) error {
	return c.handler.Transaction(fc)
}

func NewSortBy(alias, column string, direction string) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: column, Table: alias}, Desc: types.Asc.Value(direction).Bool()}
}

func NewRepository[Entity CommonEntity, Dto CommonDtoT, Relations any](handler CommonRepositoryT[Dto, Entity, Relations]) *CommonRepository[Entity, Dto, Relations] {
	return &CommonRepository[Entity, Dto, Relations]{handler}
}
