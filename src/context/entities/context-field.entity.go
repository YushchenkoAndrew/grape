package entities

import (
	"encoding/json"
	e "grape/src/common/entities"
)

type ContextFieldEntity struct {
	*e.UuidEntity
	*e.DroppableEntity

	Name    string  `gorm:"not null"`
	Value   *string `gorm:"default:null"`
	Options *string `gorm:"default:null"`

	ContextID int64 `gorm:"not null" copier:"-"`
}

func (*ContextFieldEntity) TableName() string {
	return "context_fields"
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
	return &ContextFieldEntity{
		UuidEntity:      e.NewUuidEntity(),
		DroppableEntity: e.NewDroppableEntity(),
	}
}
