package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Name      string    `gorm:"not null" json:"name,omitempty" xml:"name,omitempty" example:"metrics"`
	CronID    string    `gorm:"not null;unique" json:"cron_id" xml:"cron_id" example:"d266389ebf09e1e8a95a5b4286b504b2"`
	CronTime  string    `json:"cron_time" xml:"cron_time" example:"00 00 00 */1 * *"`
	Method    string    `json:"method" xml:"method" example:"post"`
	Path      string    `json:"path" xml:"path" example:"/ping"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
}

func NewSubscription() interfaces.Table {
	return &Subscription{}
}

func (*Subscription) TableName() string {
	return "subscription"
}

func (c *Subscription) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

type SubscribeDto struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" example:"metrics"`
	CronTime string `json:"cron_time,omitempty" xml:"cron_time,omitempty" example:"00 00 00 */1 * *"`
}

func (c *SubscribeDto) IsOK() bool {
	return c.CronTime != "" && c.Name != ""
}

type SubscribeQueryDto struct {
	ID        uint32 `form:"id,omitempty"`
	Name      string `form:"name,omitempty" example:"main"`
	ProjectID uint32 `form:"project_id,omitempty" example:"1"`

	Page  int `form:"page,omitempty" example:"1"`
	Limit int `form:"limit,omitempty" example:"10"`
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	// Link      string    `form:"link" example:"https://github.com/YushchenkoAndrew/template"`
}
