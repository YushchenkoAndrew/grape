package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type Metrics struct {
	ID            uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Name          string    `gorm:"not null" json:"name" xml:"name" example:"void-deployment-8985bd57d-k9n5g"`
	Namespace     string    `gorm:"not null" json:"namespace" xml:"namespace" example:"void-deployment-8985bd57d-k9n5g"`
	ContainerName string    `gorm:"not null" json:"container_name" xml:"container_name" example:"void"`
	CPU           int64     `gorm:"not null" json:"cpu" xml:"cpu" example:"690791"`
	Memory        int64     `gorm:"not null" json:"memory" xml:"memory" example:"690791"`
	ProjectID     uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`

	// CpuScale    int8 `gorm:"not null" json:"cpu_scale" xml:"cpu_scale" example:"3"`
	// MemoryScale int8 `gorm:"not null" json:"memory_scale" xml:"memory_scale" example:"6"`
}

func NewMetrics() interfaces.Table {
	return &Metrics{}
}

func (*Metrics) TableName() string {
	return "metrics"
}

func (c *Metrics) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

func (c *Metrics) Copy() *Metrics {
	return &Metrics{
		ID: c.ID,
		// UpdatedAt: c.UpdatedAt,
		Name: c.Name,
		// Link:      c.Link,
		ProjectID: c.ProjectID,
	}
}

func (c *Metrics) Fill(updated *Metrics) *Metrics {
	if updated.ID != 0 {
		c.ID = updated.ID
	}

	// if !updated.UpdatedAt.IsZero() {
	// 	c.UpdatedAt = updated.UpdatedAt
	// }

	if updated.Name != "" {
		c.Name = updated.Name
	}

	// if updated.Link != "" {
	// 	c.Link = updated.Link
	// }

	if updated.ProjectID != 0 {
		c.ProjectID = updated.ProjectID
	}

	return c
}

type MetricsDto struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	Link string `json:"link,omitempty" xml:"link,omitempty"`
}

func (c *MetricsDto) IsOK() bool {
	return c.Name != "" && c.Link != ""
}

type MetricsQueryDto struct {
	ID        uint32 `form:"id,omitempty"`
	Name      string `form:"name,omitempty" example:"main"`
	ProjectID uint32 `form:"project_id,omitempty" example:"1"`

	Page  int `form:"page,omitempty" example:"1"`
	Limit int `form:"limit,omitempty" example:"10"`
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	// Link      string    `form:"link" example:"https://github.com/YushchenkoAndrew/template"`
}
