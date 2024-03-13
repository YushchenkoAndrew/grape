package entities

import (
	"time"

	"github.com/google/uuid"
)

type BasicEntity struct {
	ID int64 `gorm:"primaryKey" copier:"-"`
}

func (*BasicEntity) Create() {}
func (*BasicEntity) Update() {}

func NewBasicEntity() *BasicEntity {
	return &BasicEntity{}
}

type IdEntity struct {
	*BasicEntity
	CreatedAt time.Time `gorm:"autoCreateTime;not null" copier:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" copier:"-"`
}

func (c *IdEntity) Create() { c.CreatedAt, c.UpdatedAt = time.Now(), time.Now() }
func (c *IdEntity) Update() { c.UpdatedAt = time.Now() }

func (c *IdEntity) CreatedAtISO() string { return c.CreatedAt.Format(time.RFC3339) }
func (c *IdEntity) UpdatedAtISO() string { return c.UpdatedAt.Format(time.RFC3339) }

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
