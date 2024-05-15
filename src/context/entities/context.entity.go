package entities

import (
	e "grape/src/common/entities"

	"github.com/samber/lo"
)

type ContextEntity struct {
	*e.UuidEntity
	*e.DroppableEntity

	Name string `gorm:"not null"`

	ContextableID   int64  `gorm:"not null"`
	ContextableType string `gorm:"not null"`

	ContextFields []ContextFieldEntity `gorm:"foreignKey:ContextID;references:ID" copier:"-"`
}

func (*ContextEntity) TableName() string {
	return "contexts"
}

func (c *ContextEntity) GetContextFields() []*ContextFieldEntity {
	return lo.Map(c.ContextFields, func(e ContextFieldEntity, _ int) *ContextFieldEntity { return &e })
}

func NewContextEntity() *ContextEntity {
	return &ContextEntity{
		UuidEntity:      e.NewUuidEntity(),
		DroppableEntity: e.NewDroppableEntity(),
	}
}
