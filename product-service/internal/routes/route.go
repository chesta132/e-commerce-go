package routes

import "gorm.io/gorm"

type Route struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Route {
	return &Route{db}
}
