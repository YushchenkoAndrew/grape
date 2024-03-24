package entities

import (
	"grape/src/common/config"
	e "grape/src/common/entities"
	"net/url"
	"path/filepath"
)

type AttachmentEntity struct {
	*e.UuidEntity

	Name string `gorm:"not null"`
	Home string `gorm:"not null;default:'/'"`
	Path string `gorm:"not null;default:'/'"`
	Type string `gorm:"not null"`
	Size int64  `gorm:"not null"`

	AttachableID   int64  `gorm:"not null"`
	AttachableType string `gorm:"not null"`
}

func (*AttachmentEntity) TableName() string {
	return "attachments"
}

func (c *AttachmentEntity) GetAttachment() string {
	url := url.URL{
		Scheme: "http",
		Host:   config.GetGlobalConfig().Server.Address,
		Path:   filepath.Join(config.GetGlobalConfig().Server.Prefix, c.TableName(), c.UUID),
	}
	return url.String()
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
