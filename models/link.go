package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type Link struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"name" xml:"name" example:"main"`
	Link      string    `gorm:"not null" json:"link" xml:"link" example:"https://github.com/YushchenkoAndrew/template"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	// Project   Project   `gorm:""`
}

func NewLink() interfaces.Table {
	return &Link{}
}

func (*Link) TableName() string {
	return "link"
}

func (c *Link) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

func (c *Link) Copy() *Link {
	return &Link{
		ID:        c.ID,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		Link:      c.Link,
		ProjectID: c.ProjectID,
	}
}

func (c *Link) Fill(updated *Link) *Link {
	if updated.ID != 0 {
		c.ID = updated.ID
	}

	if !updated.UpdatedAt.IsZero() {
		c.UpdatedAt = updated.UpdatedAt
	}

	if updated.Name != "" {
		c.Name = updated.Name
	}

	if updated.Link != "" {
		c.Link = updated.Link
	}

	if updated.ProjectID != 0 {
		c.ProjectID = updated.ProjectID
	}

	return c
}

type LinkDto struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Link string `json:"link" xml:"link"`
}

func (c *LinkDto) IsOK() bool {
	return c.Name != "" && c.Link != ""
}

type LinkQueryDto struct {
	ID        uint32 `form:"id"`
	Name      string `form:"name" example:"main"`
	ProjectID uint32 `form:"project_id" example:"1"`

	Page  int `form:"page" example:"1"`
	Limit int `form:"limit" example:"10"`
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	// Link      string    `form:"link" example:"https://github.com/YushchenkoAndrew/template"`
}
