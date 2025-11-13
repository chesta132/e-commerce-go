package user

import "time"

type User struct {
	ID       string `gorm:"primarykey;default:gen_random_uuid()"`
	FullName string `validate:"required"`
	Email    string `validate:"required,email" gorm:"index"`
	Password string `validate:"required"`
	Role     string `gorm:"type:user_role;default:'user'"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
