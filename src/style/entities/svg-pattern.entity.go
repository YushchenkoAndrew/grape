package entities

import (
	"encoding/json"
	"grape/src/common/entities"
	t "grape/src/style/types"
	org "grape/src/user/entities"
)

type SvgPatternEntity struct {
	*entities.UuidEntity

	OrganizationID int64                  `gorm:"not null" copier:"-"`
	Organization   org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Mode   t.ColorPaletteModeEnum `gorm:"not null,default:0"`
	Colors int                    `gorm:"not null"`

	Options string  `gorm:"not null"`
	Width   float32 `gorm:"not null"`
	Height  float32 `gorm:"not null"`

	Path string `gorm:"not null"`
}

func (*SvgPatternEntity) TableName() string {
	return "svg_patterns"
}

func (c *SvgPatternEntity) GetOptions() *t.ColorPaletteOptionsType {
	result := &t.ColorPaletteOptionsType{}
	json.Unmarshal([]byte(c.Options), result)
	return result
}

func (c *SvgPatternEntity) SetOptions(data *t.ColorPaletteOptionsType) {
	json, _ := json.Marshal(data)
	c.Options = string(json)
}

func (c *SvgPatternEntity) GetMode() string {
	return c.Mode.String()
}

func (c *SvgPatternEntity) SetMode(str string) {
	if len(str) != 0 {
		c.Mode = t.Fill.Value(str)
	}
}

func NewSvgPatternEntity() *SvgPatternEntity {
	return &SvgPatternEntity{UuidEntity: entities.NewUuidEntity()}
}
