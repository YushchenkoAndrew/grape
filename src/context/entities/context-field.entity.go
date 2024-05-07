package entities

import (
	"encoding/json"
	e "grape/src/common/entities"
)

type ContextFieldEntity struct {
	*e.UuidEntity

	Name    string  `gorm:"not null"`
	Value   *string `gorm:"default:null"`
	Order   int     `gorm:"not null;default:1" copier:"-"`
	Options *string `gorm:"default:null"`

	ContextID int64 `gorm:"not null" copier:"-"`
}

func (*ContextFieldEntity) TableName() string {
	return "context_fields"
}

func (c *ContextFieldEntity) SetOrder(order int) {
	c.Order = order
}

func (c *ContextFieldEntity) GetOrder() int {
	return c.Order
}

func (c *ContextFieldEntity) GetOptions() interface{} {
	if c.Options == nil {
		return nil
	}

	var options interface{}
	json.Unmarshal([]byte(*c.Options), options)
	return options
}

func (c *ContextFieldEntity) SetOptions(data interface{}) {
	if data == nil {
		c.Options = nil
		return
	}

	if json, err := json.Marshal(data); err != nil {
		options := string(json)
		c.Options = &options
	} else {
		c.Options = nil
	}
}

func NewContextFieldEntity() *ContextFieldEntity {
	return &ContextFieldEntity{UuidEntity: e.NewUuidEntity()}
}