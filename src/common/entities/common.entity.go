package entities

import (
	"time"

	"github.com/google/uuid"
)

type BasicEntity struct {
	ID int64 `gorm:"primaryKey" copier:"-"`
}

func (c *BasicEntity) Create() {}
func (c *BasicEntity) Update() {}

func NewBasicEntity() *BasicEntity {
	return &BasicEntity{}
}

type IdEntity struct {
	*BasicEntity
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
	return &IdEntity{BasicEntity: NewBasicEntity()}
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
