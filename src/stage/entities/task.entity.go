package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	t "grape/src/common/types"
	ctx "grape/src/context/entities"
	ln "grape/src/link/entities"
	org "grape/src/user/entities"

	"github.com/samber/lo"
)

type TaskEntity struct {
	*e.UuidEntity

	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:''"`
	Status      t.StatusEnum `gorm:"not null;default:1"`
	Order       int          `gorm:"not null;default:1" copier:"-"`

	StageID int64 `gorm:"not null" copier:"-"`

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Links       []ln.LinkEntity        `gorm:"polymorphic:Linkable" copier:"-"`
	Attachments []att.AttachmentEntity `gorm:"polymorphic:Attachable" copier:"-"`
	Contexts    []ctx.ContextEntity    `gorm:"polymorphic:Contextable" copier:"-"`
}

func (*TaskEntity) TableName() string {
	return "tasks"
}

func (c *TaskEntity) GetStatus() string {
	return c.Status.String()
}

func (c *TaskEntity) SetStatus(str string) {
	if str != "" {
		c.Status = t.Active.Value(str)
	}
}

func (c *TaskEntity) GetAttachments() []*att.AttachmentEntity {
	return lo.Map(c.Attachments, func(e att.AttachmentEntity, _ int) *att.AttachmentEntity { return &e })
}

func (c *TaskEntity) GetLinks() []*ln.LinkEntity {
	return lo.Map(c.Links, func(e ln.LinkEntity, _ int) *ln.LinkEntity { return &e })
}

func (c *TaskEntity) GetContexts() []*ctx.ContextEntity {
	return lo.Map(c.Contexts, func(e ctx.ContextEntity, _ int) *ctx.ContextEntity { return &e })
}

func NewTaskEntity() *TaskEntity {
	return &TaskEntity{UuidEntity: e.NewUuidEntity()}
}
