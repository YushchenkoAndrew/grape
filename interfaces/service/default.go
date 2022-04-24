package service

import "gorm.io/gorm"

type DefaultLogic[T any, U any] interface {
	keys(*T) []string
	isExist(*T) bool

	// Cache handler
	precache(*T, []string)
	deepcache([]T, string) interface{}
	postfilter([]T, string) []T
	recache(*T, []string, bool)

	// Query builder
	query(*U, *gorm.DB) (*gorm.DB, string)

	// Main logic
	Create(*T) error
	Read(*U) ([]T, error)
	Update(*U, *T) ([]T, error)
	Delete(*U) (int, error)
}

type Default[T any, U any] interface {
	// Main logic
	Create(*T) error
	Read(*U) ([]T, error)
	Update(*U, *T) ([]T, error)
	Delete(*U) (int, error)
}
