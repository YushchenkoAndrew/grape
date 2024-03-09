package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	t "grape/src/project/types"
	org "grape/src/user/entities"
)

type ProjectEntity struct {
	*e.UuidEntity

	Name        string              `gorm:"not null"`
	Description string              `gorm:"default:''"`
	Type        t.ProjectTypeEnum   `gorm:"not null;default:0" copier:"GetType"`
	Status      t.ProjectStatusEnum `gorm:"not null;default:0" copier:"-"`
	Footer      string              ``
	Order       int                 `gorm:"not null;default:0" copier:"-"`

	OwnerID int64          `gorm:"not null" copier:"-"`
	Owner   org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                  `gorm:"not null" copier:"-"`
	Organization   org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Links       []LinkEntity           `gorm:"-;foreignKey:ProjectID" copier:"-"`
	Attachments []att.AttachmentEntity `gorm:"-" copier:"-"`

	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

func (c *ProjectEntity) GetStatus() string {
	return c.Status.String()
}

func (c *ProjectEntity) GetType() string {
	return c.Type.String()
}

func NewProjectEntity() *ProjectEntity {
	return &ProjectEntity{UuidEntity: e.NewUuidEntity()}
}
