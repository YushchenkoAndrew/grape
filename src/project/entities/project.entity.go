package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	ln "grape/src/link/entities"
	t "grape/src/project/types"
	st "grape/src/statistic/entities"
	style "grape/src/style/entities"
	org "grape/src/user/entities"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type ProjectEntity struct {
	*e.UuidEntity

	Name        string              `gorm:"not null"`
	Description string              `gorm:"default:''"`
	Type        t.ProjectTypeEnum   `gorm:"not null;default:0"`
	Status      t.ProjectStatusEnum `gorm:"not null;default:0" copier:"-"`
	Footer      string              ``
	Order       int                 `gorm:"not null;default:0" copier:"-"`

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	ColorPaletteID int64                     `gorm:"not null" copier:"-"`
	ColorPalette   *style.ColorPaletteEntity `gorm:"foreignKey:ColorPaletteID;references:ID" copier:"-"`

	StatisticID int64               `gorm:"not null" copier:"-"`
	Statistic   *st.StatisticEntity `gorm:"foreignKey:StatisticID;references:ID" copier:"-"`

	SvgPatternID int64                   `gorm:"not null" copier:"-"`
	SvgPattern   *style.SvgPatternEntity `gorm:"foreignKey:SvgPatternID;references:ID" copier:"-"`

	Links       []ln.LinkEntity        `gorm:"foreignKey:ProjectID" copier:"-"`
	Attachments []att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

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
	if len(str) != 0 {
		c.Status = t.Active.Value(str)
	}
}

func (c *ProjectEntity) GetType() string {
	return c.Type.String()
}

func (c *ProjectEntity) SetType(str string) {
	if len(str) != 0 {
		c.Type = t.Html.Value(str)
	}
}

func (c *ProjectEntity) GetColors() []string {
	return []string(c.ColorPalette.Colors)
}

func (c *ProjectEntity) GetThumbnail() *att.AttachmentEntity {
	if result, found := lo.Find(c.Attachments, func(item att.AttachmentEntity) bool { return strings.Contains(item.Name, "thumbnail") }); found {
		return &result
	}

	return nil
}

func (c *ProjectEntity) GetPath() string {
	return filepath.Join("/", c.TableName(), c.UUID)
}

func NewProjectEntity() *ProjectEntity {
	return &ProjectEntity{UuidEntity: e.NewUuidEntity()}
}
