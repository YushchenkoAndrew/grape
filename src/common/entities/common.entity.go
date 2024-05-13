package entities

import (
	"grape/src/common/types"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type BasicEntity struct {
	ID int64 `gorm:"primaryKey" copier:"-"`
}

func (*BasicEntity) Create() {}
func (*BasicEntity) Update() {}

func (c *BasicEntity) GetID() int64 { return c.ID }

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

func (c *IdEntity) SetOrder(int)  {}
func (c *IdEntity) GetOrder() int { return 0 }

func NewIdEntity() *IdEntity {
	return &IdEntity{BasicEntity: NewBasicEntity()}
}

type UuidEntity struct {
	*IdEntity
	UUID string `gorm:"unique;not null" copier:"-"`
}

func (c *UuidEntity) TableName() string { return "" }
func (c *UuidEntity) GetPath() string   { return filepath.Join("/", c.TableName(), c.UUID) }

func (c *UuidEntity) Create() {
	c.IdEntity.Create()
	c.UUID = uuid.New().String()
}

func NewUuidEntity() *UuidEntity {
	return &UuidEntity{IdEntity: NewIdEntity()}
}

type DroppableEntity struct {
	Order int `gorm:"not null;default:1" copier:"-"`
}

func (c *DroppableEntity) GetOrder() int      { return c.Order }
func (c *DroppableEntity) SetOrder(order int) { c.Order = order }

func NewDroppableEntity() *DroppableEntity {
	return &DroppableEntity{}
}

type DeleteableEntity struct {
	Status types.StatusEnum `gorm:"not null;default:1"`
}

func (c *DeleteableEntity) GetStatus() string { return c.Status.String() }
func (c *DeleteableEntity) SetStatus(str string) {
	if str != "" {
		c.Status = types.Active.Value(str)
	}
}

func NewDeleteableEntity() *DeleteableEntity {
	return &DeleteableEntity{}
}
