package entities

import (
	"grape/src/common/entities"
	org "grape/src/user/entities"

	"github.com/lib/pq"
)

type ColorPaletteEntity struct {
	*entities.UuidEntity

	OrganizationID int64                  `gorm:"not null" copier:"-"`
	Organization   org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Colors pq.StringArray `gorm:"type:character varying[]"`
}

func (*ColorPaletteEntity) TableName() string {
	return "color_palettes"
}

func NewColorPaletteEntity() *ColorPaletteEntity {
	return &ColorPaletteEntity{UuidEntity: entities.NewUuidEntity()}
}
