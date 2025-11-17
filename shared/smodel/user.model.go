package smodel

import "time"

type User struct {
	ID       string `json:"id"`
	FullName string `validate:"required" json:"fullName"`
	Email    string `validate:"required,email" gorm:"index" json:"email"`
	Password string `validate:"required" json:"-"`
	Role     string `json:"role"`
	Address  string `json:"address"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
