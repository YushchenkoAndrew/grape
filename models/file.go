package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"name" xml:"name" example:"index.js"`
	Path      string    `json:"path" xml:"path" example:"/test"`
	Type      string    `gorm:"not null" json:"type" xml:"type" example:"js"`
	Role      string    `gorm:"not null" json:"role" xml:"role" example:"src"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	// Project   Project   `gorm:""`
}

func NewFile() interfaces.Table {
	return &File{}
}

func (*File) TableName() string {
	return "file"
}

func (c *File) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

func (c *File) Copy() *File {
	return &File{
		ID:        c.ID,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		Path:      c.Path,
		Type:      c.Type,
		Role:      c.Role,
		ProjectID: c.ProjectID,
	}
}

func (c *File) Fill(updated *File) *File {
	if updated.ID != 0 {
		c.ID = updated.ID
	}

	if !updated.UpdatedAt.IsZero() {
		c.UpdatedAt = updated.UpdatedAt
	}

	if updated.Name != "" {
		c.Name = updated.Name
	}

	if updated.Path != "" {
		c.Path = updated.Path
	}

	if updated.Type != "" {
		c.Type = updated.Type
	}

	if updated.Role != "" {
		c.Role = updated.Role
	}

	if updated.ProjectID != 0 {
		c.ProjectID = updated.ProjectID
	}

	return c
}

type FileDto struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Path string `json:"path" xml:"path"`
	Type string `json:"type" xml:"type"`
	Role string `json:"role" xml:"role"`
}

func (c *FileDto) IsOK() bool {
	return c.Name == "" || c.Type == ""
}

type FileQueryDto struct {
	ID        uint32 `form:"id"`
	Name      string `form:"name" example:"main"`
	Role      string `form:"role" example:"src"`
	Path      string `form:"path" example:"/test"`
	ProjectID uint32 `form:"project_id" example:"1"`

	Page  int `form:"page" example:"1"`
	Limit int `form:"limit" example:"10"`
}
