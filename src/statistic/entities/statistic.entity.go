package entities

import (
	e "grape/src/common/entities"
)

type StatisticEntity struct {
	*e.UuidEntity

	Views  int `gorm:"default:0"`
	Clicks int `gorm:"default:0"`
	Media  int `gorm:"default:0"`
}

func (*StatisticEntity) TableName() string {
	return "statistics"
}

func NewStatisticEntity() *StatisticEntity {
	return &StatisticEntity{UuidEntity: e.NewUuidEntity()}
}
