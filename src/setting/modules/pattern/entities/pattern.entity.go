package entities

import (
	"encoding/json"
	"grape/src/common/entities"
	t "grape/src/setting/modules/pattern/types"
	org "grape/src/user/entities"
)

type PatternEntity struct {
	*entities.UuidEntity

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Mode   t.PatternColorModeEnum `gorm:"not null,default:1"`
	Colors int                    `gorm:"not null"`

	Options string  `gorm:"not null"`
	Width   float32 `gorm:"not null"`
	Height  float32 `gorm:"not null"`

	Path string `gorm:"not null"`
}

func (*PatternEntity) TableName() string {
	return "patterns"
}

func (c *PatternEntity) GetOptions() *t.PatternOptionsType {
	result := &t.PatternOptionsType{}
	json.Unmarshal([]byte(c.Options), result)
	return result
}

func (c *PatternEntity) SetOptions(data *t.PatternOptionsType) {
	json, _ := json.Marshal(data)
	c.Options = string(json)
}

func (c *PatternEntity) GetMode() string {
	return c.Mode.String()
}

func (c *PatternEntity) SetMode(str string) {
	if len(str) != 0 {
		c.Mode = t.Fill.Value(str)
	}
}

func NewPatternEntity() *PatternEntity {
	return &PatternEntity{UuidEntity: entities.NewUuidEntity()}
}
