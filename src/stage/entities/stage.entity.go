package entities

import (
	e "grape/src/common/entities"
	t "grape/src/common/types"
	org "grape/src/user/entities"

	"github.com/samber/lo"
)

type StageEntity struct {
	*e.UuidEntity

	Name   string       `gorm:"not null"`
	Status t.StatusEnum `gorm:"not null;default:1"`
	Order  int          `gorm:"not null;default:1" copier:"-"`

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Tasks []TaskEntity `gorm:"foreignKey:StageID;references:ID" copier:"-"`
}

func (*StageEntity) TableName() string {
	return "stages"
}

func (c *StageEntity) SetOrder(order int) {
	c.Order = order
}

func (c *StageEntity) GetOrder() int {
	return c.Order
}

func (c *StageEntity) GetStatus() string {
	return c.Status.String()
}

func (c *StageEntity) SetStatus(str string) {
	if str != "" {
		c.Status = t.Active.Value(str)
	}
}

func (c *StageEntity) GetTasks() []*TaskEntity {
	return lo.Map(c.Tasks, func(e TaskEntity, _ int) *TaskEntity { return &e })
}

func NewStageEntity() *StageEntity {
	return &StageEntity{UuidEntity: e.NewUuidEntity()}
}
