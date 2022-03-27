package models

import (
	"api/interfaces"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Link struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
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

func (c *Link) Migrate(db *gorm.DB, forced bool) {
	if forced {
		db.Migrator().DropTable(c)
	}
	db.AutoMigrate(c)
}

func (c *Link) Redis(db *gorm.DB, client *redis.Client) error {
	var value int64
	db.Model(c).Count(&value)

	if err := client.Set(context.Background(), "nLINK", value, 0).Err(); err != nil {
		return fmt.Errorf("[Redis] Error happed while setting value to Cache: %v", err)
	}

	return nil
}

func (c *Link) ToModel(dto *LinkDto) {
	c.Name = dto.Name
	c.Link = dto.Link
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
