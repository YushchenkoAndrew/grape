package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	t "grape/src/project/types"
	org "grape/src/user/entities"
)

type ProjectEntity struct {
	e.UuidEntity

	Name        string              `gorm:"not null"`
	Description string              `gorm:"default:''"`
	Type        t.ProjectTypeEnum   `gorm:"not null;default:0"`
	Status      t.ProjectStatusEnum `gorm:"not null;default:0"`
	Footer      string              ``

	OwnerID int64          `gorm:"not null"`
	Owner   org.UserEntity `gorm:"foreignKey:OwnerID;references:ID"`

	OrganizationID int64                  `gorm:"not null"`
	Organization   org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID"`

	Links       []LinkEntity           `gorm:"-;foreignKey:ProjectID"`
	Attachments []att.AttachmentEntity `gorm:"-"`

	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

func (c *ProjectEntity) GetStatus() string {
	switch c.Status {
	case t.Active:
		return "active"

	case t.Inactive:
		return "inactive"
	}

	return ""
}

func (c *ProjectEntity) GetType() string {
	switch c.Type {
	case t.Html:
		return "html"

	case t.Markdown:
		return "markdown"

	case t.Link:
		return "link"

	case t.K3s:
		return "k3s"
	}

	return ""
}
