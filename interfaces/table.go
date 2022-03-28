package interfaces

import (
	"gorm.io/gorm"
)

type Table interface {
	Migrate(db *gorm.DB, forced bool) error
	// Redis(db *gorm.DB, client *redis.Client) error
}
