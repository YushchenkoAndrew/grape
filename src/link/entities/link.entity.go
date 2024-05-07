package entities

import (
	e "grape/src/common/entities"
)

type LinkEntity struct {
	*e.UuidEntity

	Name  string `gorm:"not null"`
	Link  string `gorm:"not null"`
	Order int    `gorm:"not null;default:1" copier:"-"`

	LinkableID   int64  `gorm:"not null" copier:"-"`
	LinkableType string `gorm:"not null" copier:"-"`
}

func (*LinkEntity) TableName() string {
	return "links"
}

func (c *LinkEntity) SetOrder(order int) {
	c.Order = order
}

func (c *LinkEntity) GetOrder() int {
	return c.Order
}

func NewLinkEntity() *LinkEntity {
	return &LinkEntity{UuidEntity: e.NewUuidEntity()}
}
