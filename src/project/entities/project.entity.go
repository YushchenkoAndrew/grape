package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	ln "grape/src/link/entities"
	t "grape/src/project/types"
	palette "grape/src/setting/modules/palette/entities"
	pattern "grape/src/setting/modules/pattern/entities"
	st "grape/src/statistic/entities"
	org "grape/src/user/entities"
	"path/filepath"
)

type ProjectEntity struct {
	*e.UuidEntity

	Name        string              `gorm:"not null"`
	Description string              `gorm:"default:''"`
	Type        t.ProjectTypeEnum   `gorm:"not null;default:1"`
	Status      t.ProjectStatusEnum `gorm:"not null;default:1"`
	Footer      string              ``
	Order       int                 `gorm:"not null;default:0" copier:"-"`

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	PaletteID int64                  `gorm:"not null" copier:"-"`
	Palette   *palette.PaletteEntity `gorm:"foreignKey:PaletteID;references:ID" copier:"-"`

	StatisticID int64               `gorm:"not null" copier:"-"`
	Statistic   *st.StatisticEntity `gorm:"foreignKey:StatisticID;references:ID" copier:"-"`

	PatternID int64                  `gorm:"not null" copier:"-"`
	Pattern   *pattern.PatternEntity `gorm:"foreignKey:PatternID;references:ID" copier:"-"`

	Links       []ln.LinkEntity        `gorm:"polymorphic:Linkable" copier:"-"`
	Attachments []att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

	Thumbnail *att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

func (c *ProjectEntity) GetStatus() string {
	return c.Status.String()
}

func (c *ProjectEntity) SetStatus(str string) {
	if str != "" {
		c.Status = t.Active.Value(str)
	}
}

func (c *ProjectEntity) GetType() string {
	return c.Type.String()
}

func (c *ProjectEntity) SetType(str string) {
	if str != "" {
		c.Type = t.Html.Value(str)
	}
}

func (c *ProjectEntity) GetPath() string {
	return filepath.Join("/", c.TableName(), c.UUID)
}

func NewProjectEntity() *ProjectEntity {
	return &ProjectEntity{UuidEntity: e.NewUuidEntity()}
}
