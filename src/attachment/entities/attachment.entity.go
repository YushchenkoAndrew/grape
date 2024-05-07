package entities

import (
	"fmt"
	e "grape/src/common/entities"
	"path/filepath"
)

type AttachmentEntity struct {
	*e.UuidEntity

	Name    string `gorm:"not null"`
	Home    string `gorm:"not null;default:'/'"`
	Path    string `gorm:"not null;default:'/'"`
	Type    string `gorm:"not null"`
	Size    int64  `gorm:"not null"`
	Order   int    `gorm:"not null;default:1" copier:"-"`
	Preview bool   `gorm:"not null"`

	AttachableID   int64  `gorm:"not null"`
	AttachableType string `gorm:"not null"`
}

func (*AttachmentEntity) TableName() string {
	return "attachments"
}

func (c *AttachmentEntity) SetOrder(order int) {
	c.Order = order
}

func (c *AttachmentEntity) GetOrder() int {
	return c.Order
}

func (c *AttachmentEntity) GetAttachment() string {
	return fmt.Sprintf("/attachments/%s", c.UUID)
}

func (c *AttachmentEntity) GetPath() string {
	return filepath.Join(c.Home, c.Path)
}

func (c *AttachmentEntity) GetFile() string {
	return filepath.Join(c.GetPath(), c.Name)
}

func NewAttachmentEntity() *AttachmentEntity {
	return &AttachmentEntity{UuidEntity: e.NewUuidEntity()}
}
