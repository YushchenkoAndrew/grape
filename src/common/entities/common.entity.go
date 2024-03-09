package entities

import (
	"time"

	"github.com/google/uuid"
)

type IdEntity struct {
	ID        int64     `gorm:"primaryKey" copier:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" copier:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" copier:"-"`
}

func (c *IdEntity) Create() {
	c.CreatedAt, c.UpdatedAt = time.Now(), time.Now()
}

func (c *IdEntity) Update() {
	c.UpdatedAt = time.Now()
}

func NewIdEntity() *IdEntity {
	return &IdEntity{}
}

type UuidEntity struct {
	*IdEntity
	UUID string `gorm:"unique;not null" copier:"-"`
}

func (c *UuidEntity) Create() {
	c.IdEntity.Create()
	c.UUID = uuid.New().String()
}

func NewUuidEntity() *UuidEntity {
	return &UuidEntity{IdEntity: NewIdEntity()}
}
