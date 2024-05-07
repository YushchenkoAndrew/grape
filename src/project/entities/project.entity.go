package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	"grape/src/common/types"
	ln "grape/src/link/entities"
	t "grape/src/project/types"
	st "grape/src/statistic/entities"
	org "grape/src/user/entities"
	"path/filepath"

	"github.com/samber/lo"
)

type ProjectEntity struct {
	*e.UuidEntity

	Name        string            `gorm:"not null"`
	Description string            `gorm:"default:''"`
	Type        t.ProjectTypeEnum `gorm:"not null;default:1"`
	Status      types.StatusEnum  `gorm:"not null;default:1"`
	Footer      string            ``
	Order       int               `gorm:"not null;default:1" copier:"-"`

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	StatisticID int64               `gorm:"not null" copier:"-"`
	Statistic   *st.StatisticEntity `gorm:"foreignKey:StatisticID;references:ID" copier:"-"`

	Links       []ln.LinkEntity        `gorm:"polymorphic:Linkable" copier:"-"`
	Attachments []att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

	Redirect  *ln.LinkEntity        `gorm:"polymorphic:Linkable" copier:"-"`
	Thumbnail *att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`

	// Metrics      []Metrics      `gorm:"foreignKey:ProjectID" json:"metrics" xml:"metrics"`
	// Subscription []Subscription `gorm:"foreignKey:ProjectID" json:"subscription" xml:"subscription"`
}

func (*ProjectEntity) TableName() string {
	return "projects"
}

func (c *ProjectEntity) SetOrder(order int) {
	c.Order = order
}

func (c *ProjectEntity) GetOrder() int {
	return c.Order
}

func (c *ProjectEntity) GetStatus() string {
	return c.Status.String()
}

func (c *ProjectEntity) SetStatus(str string) {
	if str != "" {
		c.Status = types.Active.Value(str)
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

func (c *ProjectEntity) GetAttachments() []*att.AttachmentEntity {
	return lo.Map(c.Attachments, func(e att.AttachmentEntity, _ int) *att.AttachmentEntity { return &e })
}

func (c *ProjectEntity) GetLinks() []*ln.LinkEntity {
	return lo.Map(c.Links, func(e ln.LinkEntity, _ int) *ln.LinkEntity { return &e })
}

func NewProjectEntity() *ProjectEntity {
	return &ProjectEntity{UuidEntity: e.NewUuidEntity()}
}
