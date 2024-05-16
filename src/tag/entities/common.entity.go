package entities

import "github.com/samber/lo"

type TaggableEntity struct {
	Tags []TagEntity `gorm:"polymorphic:Taggable" copier:"-"`
}

func (c *TaggableEntity) GetTags() []*TagEntity {
	return lo.Map(c.Tags, func(e TagEntity, _ int) *TagEntity { return &e })
}

func NewTaggableEntity() *TaggableEntity {
	return &TaggableEntity{}
}

type TaggableT interface {
	GetID() int64
	GetTags() []*TagEntity
}
