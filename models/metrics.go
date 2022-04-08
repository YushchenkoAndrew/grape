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
		ID:            c.ID,
		CreatedAt:     c.CreatedAt,
		Name:          c.Name,
		Namespace:     c.Namespace,
		ContainerName: c.ContainerName,
		CPU:           c.CPU,
		Memory:        c.Memory,
		ProjectID:     c.ProjectID,
	}
}

func (c *Metrics) Fill(updated *Metrics) *Metrics {
	if updated.ID != 0 {
		c.ID = updated.ID
	}

	if !updated.CreatedAt.IsZero() {
		c.CreatedAt = updated.CreatedAt
	}

	if updated.Name != "" {
		c.Name = updated.Name
	}

	if updated.Namespace != "" {
		c.Namespace = updated.Namespace
	}

	if updated.ContainerName != "" {
		c.ContainerName = updated.ContainerName
	}

	if updated.CPU != 0 {
		c.CPU = updated.CPU
	}

	if updated.Memory != 0 {
		c.Memory = updated.Memory
	}

	if updated.ProjectID != 0 {
		c.ProjectID = updated.ProjectID
	}

	return c
}

type MetricsDto struct {
	CPU    int64 `json:"cpu,omitempty" xml:"cpu,omitempty" example:"690791"`
	Memory int64 `json:"memory,omitempty" xml:"memory,omitempty" example:"690791"`
}

func (c *MetricsDto) IsOK() bool {
	return c.CPU > 0 && c.Memory > 0
}

type MetricsQueryDto struct {
	ID            uint32 `form:"id,omitempty"`
	Name          string `form:"name,omitempty" example:"main"`
	Namespace     string `form:"namespace,omitempty" example:"void-deployment-8985bd57d-k9n5g"`
	ContainerName string `form:"container_name,omitempty" example:"void"`
	ProjectID     uint32 `form:"project_id,omitempty" example:"1"`

	CreatedTo   time.Time `form:"created_to,omitempty" time_format:"2006-01-02" example:"2021-08-06"`
	CreatedFrom time.Time `form:"created_from,omitempty" time_format:"2006-01-02" example:"2021-08-06"`

	Page  int `form:"page,omitempty" example:"1"`
	Limit int `form:"limit,omitempty" example:"10"`
}
