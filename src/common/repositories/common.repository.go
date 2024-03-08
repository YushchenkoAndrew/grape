package repositories

import "gorm.io/gorm"

type CommonRepository[Entity any, Dto any, Relations any] struct {
	// db *gorm.DB
}

func (c *CommonRepository[Entity, Dto, Relations]) Build(dto *Dto, relations ...Relations) *gorm.DB {
	return nil
}

func (c *CommonRepository[Entity, Dto, Relations]) GetOne(dto *Dto, relations ...Relations) *Entity {
	var result *Entity
	c.Build(dto, relations...).Limit(1).Find(&result)
	return result
}

func (c *CommonRepository[Entity, Dto, Relations]) GetAll(dto *Dto, relations ...Relations) *[]Entity {
	var result *[]Entity
	c.Build(dto, relations...).Find(&result)
	return result
}

// func (c *CommonRepository[Entity, Dto, Relations]) GetAllPage(dto *Dto, relations ...Relations) *[]Entity {
// 	return nil
// }
