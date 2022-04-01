package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Name      string    `gorm:"not null;unique" json:"name" xml:"name" example:"CodeRain"`
	Title     string    `gorm:"not null" json:"title" xml:"title" example:"Code Rain"`
	Flag      string    `json:"flag" xml:"flag" example:"js"`
	Desc      string    `json:"desc" xml:"desc" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Note      string    `json:"note" xml:"note" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`

	Files        []File         `gorm:"foreignKey:ProjectID" json:"files" xml:"files"`
	Links        []Link         `gorm:"foreignKey:ProjectID" json:"links" xml:"links"`
	Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func NewProject() interfaces.Table {
	return &Project{}
}

func (*Project) TableName() string {
	return "project"
}

func (c *Project) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

func (c *Project) Copy() *Project {
	return &Project{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		Name:      c.Name,
		Title:     c.Title,
		Flag:      c.Flag,
		Desc:      c.Desc,
		Note:      c.Note,
	}
}

func (c *Project) Fill(updated *Project) *Project {
	if updated.ID != 0 {
		c.ID = updated.ID
	}

	if !updated.CreatedAt.IsZero() {
		c.CreatedAt = updated.CreatedAt
	}

	if updated.Name != "" {
		c.Name = updated.Name
	}

	if updated.Title != "" {
		c.Title = updated.Title
	}

	if updated.Flag != "" {
		c.Flag = updated.Flag
	}

	if updated.Desc != "" {
		c.Desc = updated.Desc
	}

	if updated.Note != "" {
		c.Note = updated.Note
	}

	return c
}

type ProjectDto struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name  string `json:"name,omitempty" xml:"name,omitempty"`
	Title string `json:"title,omitempty" xml:"title,omitempty"`
	Flag  string `json:"flag,omitempty" xml:"flag,omitempty"`
	Desc  string `json:"desc,omitempty" xml:"desc,omitempty"`
	Note  string `json:"note,omitempty" xml:"note,omitempty"`

	Links []LinkDto `json:"links,omitempty" xml:"links,omitempty"`
	Files []FileDto `json:"files,omitempty" xml:"files,omitempty"`
}

func (c *ProjectDto) IsOK() bool {
	return c.Name != "" && c.Title != "" && c.Flag != "" && c.Desc != ""
}

type ProjectQueryDto struct {
	ID   uint32 `form:"id,omitempty"`
	Name string `form:"name,omitempty" example:"main"`
	Flag string `form:"flag,omitempty" example:"js"`

	CreatedTo   time.Time `form:"created_to,omitempty" time_format:"2006-01-02" example:"2021-08-06"`
	CreatedFrom time.Time `form:"created_from,omitempty" time_format:"2006-01-02" example:"2021-08-06"`

	Page  int `form:"page,omitempty" example:"1"`
	Limit int `form:"limit,omitempty" example:"10"`

	Link struct {
		ID   uint32 `form:"link[id],omitempty"`
		Name string `form:"link[name],omitempty" example:"main"`

		Page  int `form:"link[page],omitempty" example:"1"`
		Limit int `form:"link[limit],omitempty" example:"10"`
	}

	File struct {
		ID   uint32 `form:"file[id],omitempty"`
		Name string `form:"file[name],omitempty" example:"main"`
		Role string `form:"file[role],omitempty" example:"src"`
		Path string `form:"file[path],omitempty" example:"/test"`

		Page  int `form:"file[page],omitempty" example:"1"`
		Limit int `form:"file[limit],omitempty" example:"10"`
	}

	Subscription struct {
		ID     uint32 `form:"subscription[id],omitempty"`
		Name   string `form:"subscription[name],omitempty" example:"main"`
		CronID string `form:"subscription[cron_id],omitempty" example:"main"`

		Page  int `form:"subscription[page],omitempty" example:"1"`
		Limit int `form:"subscription[limit],omitempty" example:"10"`
	}
}
