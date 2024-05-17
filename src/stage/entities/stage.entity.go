package entities

import (
	e "grape/src/common/entities"
	org "grape/src/user/entities"

	"github.com/samber/lo"
)

type StageEntity struct {
	*e.UuidEntity
	*e.DroppableEntity
	*e.DeleteableEntity

	Name string `gorm:"not null"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Tasks []TaskEntity `gorm:"foreignKey:StageID;references:ID" copier:"-"`
}

func (*StageEntity) TableName() string {
	return "stages"
}

func (c *StageEntity) GetTasks() []*TaskEntity {
	return lo.Map(c.Tasks, func(e TaskEntity, _ int) *TaskEntity { return &e })
}

func NewStageEntity() *StageEntity {
	return &StageEntity{
		UuidEntity:       e.NewUuidEntity(),
		DroppableEntity:  e.NewDroppableEntity(),
		DeleteableEntity: e.NewDeleteableEntity(),
	}
}
