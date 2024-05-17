package entities

import "github.com/samber/lo"

type LinkableEntity struct {
	Links []LinkEntity `gorm:"polymorphic:Linkable" copier:"-"`
}

func (c *LinkableEntity) GetLinks() []*LinkEntity {
	return lo.Map(c.Links, func(e LinkEntity, _ int) *LinkEntity { return &e })
}

func NewLinkableEntity() *LinkableEntity {
	return &LinkableEntity{}
}

type LinkableT interface {
	GetID() int64
	GetLinks() []*LinkEntity
}
