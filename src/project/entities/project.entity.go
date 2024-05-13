package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	ln "grape/src/link/entities"
	t "grape/src/project/types"
	st "grape/src/statistic/entities"
	org "grape/src/user/entities"
)

type ProjectEntity struct {
	*e.UuidEntity
	*e.DroppableEntity
	*e.DeleteableEntity

	Name        string            `gorm:"not null"`
	Description string            `gorm:"default:''"`
	Type        t.ProjectTypeEnum `gorm:"not null;default:1"`
	Footer      string            ``

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	StatisticID int64               `gorm:"not null" copier:"-"`
	Statistic   *st.StatisticEntity `gorm:"foreignKey:StatisticID;references:ID" copier:"-"`

	Redirect  *ln.LinkEntity        `gorm:"polymorphic:Linkable" copier:"-"`
	Thumbnail *att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`

	*ln.LinkableEntity
	*att.AttachableEntity
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

func (c *ProjectEntity) GetType() string {
	return c.Type.String()
}

func (c *ProjectEntity) SetType(str string) {
	if str != "" {
		c.Type = t.Html.Value(str)
	}
}

func NewProjectEntity() *ProjectEntity {
	return &ProjectEntity{
		UuidEntity:       e.NewUuidEntity(),
		DroppableEntity:  e.NewDroppableEntity(),
		DeleteableEntity: e.NewDeleteableEntity(),

		LinkableEntity:   ln.NewLinkableEntity(),
		AttachableEntity: att.NewAttachableEntity(),
	}
}
