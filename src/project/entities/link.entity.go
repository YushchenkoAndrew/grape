package entities

import (
	e "grape/src/common/entities"
)

type LinkEntity struct {
	e.UuidEntity

	Name      string `gorm:"not null"`
	Link      string `gorm:"not null"`
	ProjectID int64  `gorm:"foreignKey:ProjectID;not null"`
}

func (*LinkEntity) TableName() string {
	return "links"
}

// type LinkDto struct {
// 	// ID        uint32    `json:"id" xml:"id"`
// 	Name string `json:"name,omitempty" xml:"name,omitempty"`
// 	Link string `json:"link,omitempty" xml:"link,omitempty"`
// }

// func (c *LinkDto) IsOK() bool {
// 	return c.Name != "" && c.Link != ""
// }

// type LinkQueryDto struct {
// 	ID        uint32 `form:"id,omitempty"`
// 	Name      string `form:"name,omitempty" example:"main"`
// 	ProjectID uint32 `form:"project_id,omitempty" example:"1"`

// 	Page  int `form:"page,omitempty,default=-1" example:"1"`
// 	Limit int `form:"limit,omitempty" example:"10"`
// 	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
// 	// Link      string    `form:"link" example:"https://github.com/YushchenkoAndrew/template"`
// }
