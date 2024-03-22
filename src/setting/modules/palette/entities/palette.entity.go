package entities

import (
	"grape/src/common/entities"
	org "grape/src/user/entities"

	"github.com/lib/pq"
)

type PaletteEntity struct {
	*entities.UuidEntity

	OrganizationID int64                   `gorm:"not null" copier:"-"`
	Organization   *org.OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID" copier:"-"`

	Colors pq.StringArray `gorm:"type:character varying[]"`
}

func (*PaletteEntity) TableName() string {
	return "palettes"
}

func NewPaletteEntity() *PaletteEntity {
	return &PaletteEntity{UuidEntity: entities.NewUuidEntity()}
}
