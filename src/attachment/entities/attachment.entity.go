package entities

import (
	t "grape/src/attachment/types"
	e "grape/src/common/entities"
)

type AttachmentEntity struct {
	e.UuidEntity

	Name string `gorm:"not null"`
	Path string `gorm:"not null;default:'/'"`
	Type string `gorm:"not null"`

	AttachableID   int64                `gorm:"not null"`
	AttachableType t.AttachableTypeEnum `gorm:"not null"`
}

func (*AttachmentEntity) TableName() string {
	return "attachments"
}
