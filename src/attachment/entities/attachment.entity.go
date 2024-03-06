package entities

import (
	e "grape/src/common/entities"
)

type AttachmentEntity struct {
	e.UuidEntity

	Name      string `gorm:"not null" json:"name" xml:"name" example:"index.js"`
	Path      string `json:"path" xml:"path" example:"/test"`
	Type      string `gorm:"not null" json:"type" xml:"type" example:"js"`
	Role      string `gorm:"not null" json:"role" xml:"role" example:"src"`
	ProjectID uint32 `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	// Project   Project   `gorm:""`
}

func (*AttachmentEntity) TableName() string {
	return "attachments"
}

// type FileDto struct {
// 	// ID        uint32    `json:"id" xml:"id"`
// 	Name string `json:"name,omitempty" xml:"name,omitempty"`
// 	Path string `json:"path,omitempty" xml:"path,omitempty"`
// 	Type string `json:"type,omitempty" xml:"type,omitempty"`
// 	Role string `json:"role,omitempty" xml:"role,omitempty"`
// }

// func (c *FileDto) IsOK() bool {
// 	return c.Name != "" || c.Type != ""
// }

// type FileQueryDto struct {
// 	ID        uint32 `form:"id,omitempty"`
// 	Name      string `form:"name,omitempty" example:"main"`
// 	Role      string `form:"role,omitempty" example:"src"`
// 	Path      string `form:"path,omitempty" example:"/test"`
// 	ProjectID uint32 `form:"project_id,omitempty" example:"1"`

// 	Page  int `form:"page,omitempty,default=-1" example:"1"`
// 	Limit int `form:"limit,omitempty" example:"10"`
// }
