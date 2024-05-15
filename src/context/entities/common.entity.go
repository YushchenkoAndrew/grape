package entities

import "github.com/samber/lo"

type ContextableEntity struct {
	Contexts []ContextEntity `gorm:"polymorphic:Contextable" copier:"-"`
}

func (c *ContextableEntity) GetContexts() []*ContextEntity {
	return lo.Map(c.Contexts, func(e ContextEntity, _ int) *ContextEntity { return &e })
}

func NewContextableEntity() *ContextableEntity {
	return &ContextableEntity{}
}

type ContextableT interface {
	GetID() int64
	GetContexts() []*ContextEntity
}
