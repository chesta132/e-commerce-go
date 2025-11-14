package user

import "time"

type User struct {
	ID       string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	FullName string `validate:"required" json:"fullName"`
	Email    string `validate:"required,email" gorm:"index" json:"email"`
	Password string `validate:"required" json:"-"`
	Role     string `gorm:"type:user_role;default:'user'" json:"role"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
