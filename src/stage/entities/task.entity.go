package entities

import (
	att "grape/src/attachment/entities"
	e "grape/src/common/entities"
	ctx "grape/src/context/entities"
	ln "grape/src/link/entities"
	org "grape/src/user/entities"
)

type TaskEntity struct {
	*e.UuidEntity
	*e.DroppableEntity
	*e.DeleteableEntity

	Name        string `gorm:"not null"`
	Description string `gorm:"default:''"`

	StageID int64 `gorm:"not null" copier:"-"`

	OwnerID int64           `gorm:"not null" copier:"-"`
	Owner   *org.UserEntity `gorm:"foreignKey:OwnerID;references:ID" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	*ln.LinkableEntity
	*att.AttachableEntity
	*ctx.ContextableEntity
}

func (*TaskEntity) TableName() string {
	return "tasks"
}

func NewTaskEntity() *TaskEntity {
	return &TaskEntity{
		UuidEntity:       e.NewUuidEntity(),
		DroppableEntity:  e.NewDroppableEntity(),
		DeleteableEntity: e.NewDeleteableEntity(),

		LinkableEntity:    ln.NewLinkableEntity(),
		AttachableEntity:  att.NewAttachableEntity(),
		ContextableEntity: ctx.NewContextableEntity(),
	}
}
