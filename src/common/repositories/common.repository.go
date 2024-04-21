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

	Build(*gorm.DB, Dto, ...Relations) *gorm.DB
	Create(*gorm.DB, Dto, interface{}, Entity) *gorm.DB
	Update(*gorm.DB, Dto, interface{}, Entity) *gorm.DB
	Delete(*gorm.DB, Dto, []Entity) *gorm.DB
	Reorder(*gorm.DB, Entity, int) error
}

type CommonRepository[Entity CommonEntity, Dto CommonDtoT, Relations any] struct {
	db      *gorm.DB
	handler CommonRepositoryT[Dto, Entity, Relations]
}

func (c *CommonRepository[Entity, Dto, Relations]) connection(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return c.db
}

func (c *CommonRepository[Entity, Dto, Relations]) TableName() string {
	return c.handler.Model().TableName()
}

func (c *CommonRepository[Entity, Dto, Relations]) Transaction(fc func(tx *gorm.DB) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		return fc(tx)
	})
}

func (c *CommonRepository[Entity, Dto, Relations]) GetOne(dto Dto, relations ...Relations) (Entity, error) {
	var result []Entity
	if tx := c.handler.Build(c.db, dto, relations...).Limit(1).Find(&result); tx.Error != nil {
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
	if tx := c.handler.Build(c.db, dto, relations...).Find(&result); tx.Error != nil {
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

	if tx := c.handler.Build(c.db, dto, relations...).Count(&cnt); tx.Error != nil {
		return 0, nil, tx.Error
	}

	if tx := c.handler.Build(c.db, dto, relations...).Offset(dto.Offset()).Limit(dto.Limit()).Find(&result); tx.Error != nil {
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

func (c *CommonRepository[Entity, Dto, Relations]) Create(db *gorm.DB, dto Dto, body interface{}) (Entity, error) {
	entity := c.handler.Model()
	copier.Copy(&entity, body)
	entity.Create()

	if tx := c.handler.Create(c.connection(db), dto, body, entity); tx.Error != nil {
		var e Entity
		return e, tx.Error
	}

	return entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Update(db *gorm.DB, dto Dto, body interface{}, entity Entity) (Entity, error) {
	copier.CopyWithOption(&entity, body, copier.Option{IgnoreEmpty: true})
	entity.Update()

	if tx := c.handler.Update(c.connection(db), dto, body, entity); tx.Error != nil {
		var e Entity
		return e, tx.Error
	}

	return entity, nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Delete(db *gorm.DB, dto Dto, entity Entity) error {
	if tx := c.handler.Delete(c.connection(db), dto, []Entity{entity}); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *CommonRepository[Entity, Dto, Relations]) DeleteAll(db *gorm.DB, dto Dto, entities []Entity) error {
	if len(entities) == 0 {
		return nil
	}

	if tx := c.handler.Delete(c.connection(db), dto, entities); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *CommonRepository[Entity, Dto, Relations]) Reorder(db *gorm.DB, entity Entity, position int) error {
	if err := c.handler.Reorder(c.connection(db), entity, position); err != nil {
		return err
	}

	return nil
}

func NewSortBy(alias, column string, direction string) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: column, Table: alias}, Desc: types.Asc.Value(direction).Bool()}
}

func NewRepository[Entity CommonEntity, Dto CommonDtoT, Relations any](db *gorm.DB, handler CommonRepositoryT[Dto, Entity, Relations]) *CommonRepository[Entity, Dto, Relations] {
	return &CommonRepository[Entity, Dto, Relations]{db, handler}
}
