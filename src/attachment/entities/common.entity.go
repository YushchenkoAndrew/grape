package entities

import "github.com/samber/lo"

type AttachableEntity struct {
	Attachments []AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`
}

func (c *AttachableEntity) GetAttachments() []*AttachmentEntity {
	return lo.Map(c.Attachments, func(e AttachmentEntity, _ int) *AttachmentEntity { return &e })
}

func NewAttachableEntity() *AttachableEntity {
	return &AttachableEntity{}
}

type AttachableT interface {
	GetID() int64
	GetPath() string
	GetAttachments() []*AttachmentEntity
}
