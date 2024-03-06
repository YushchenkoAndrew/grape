package entities

import (
	e "grape/src/common/entities"
)

type ProjectEntity struct {
	*e.UuidEntity

	// Title       string `gorm:"not null" json:"title" xml:"title" example:"Code Rain"`

	Name        string `gorm:"not null" json:"name" xml:"name" example:"Code Rain"`
	Description string `json:"desc" xml:"desc" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Type        string `json:"flag" xml:"flag" example:"js"`
	Footer      string `json:"note" xml:"note" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`

	// Files        []File         `gorm:"foreignKey:ProjectID" json:"files" xml:"files"`
	// Links        []Link         `gorm:"foreignKey:ProjectID" json:"links" xml:"links"`
	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

// type ProjectDto struct {
// 	// ID        uint32    `json:"id" xml:"id"`
// 	Name  string `json:"name,omitempty" xml:"name,omitempty"`
// 	Title string `json:"title,omitempty" xml:"title,omitempty"`
// 	Flag  string `json:"flag,omitempty" xml:"flag,omitempty"`
// 	Desc  string `json:"desc,omitempty" xml:"desc,omitempty"`
// 	Note  string `json:"note,omitempty" xml:"note,omitempty"`

// 	Links []LinkDto `json:"links,omitempty" xml:"links,omitempty"`
// 	Files []FileDto `json:"files,omitempty" xml:"files,omitempty"`
// }

// func (c *ProjectDto) IsOK() bool {
// 	return c.Name != "" && c.Title != "" && c.Flag != "" && c.Desc != ""
// }

// func (c *ProjectQueryDto) IsOK(model *Project) bool {
// 	return (c.Name == "" || c.Name == model.Name) && (c.Flag == "" || c.Flag == model.Flag)
// }
