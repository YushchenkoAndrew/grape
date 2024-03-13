package entities

import (
	e "grape/src/common/entities"
)

type AttachmentEntity struct {
	*e.UuidEntity

	Name string `gorm:"not null"`
	Path string `gorm:"not null;default:'/'"`
	Type string `gorm:"not null"`
	Size int64  `gorm:"not null"`

	AttachableID   int64  `gorm:"not null"`
	AttachableType string `gorm:"not null"`
}

func (*AttachmentEntity) TableName() string {
	return "attachments"
}

func (c *AttachmentEntity) GetPath() string {
	return "CHANGE ME" + c.Path
}

func NewAttachmentEntity() *AttachmentEntity {
	return &AttachmentEntity{UuidEntity: e.NewUuidEntity()}
}
