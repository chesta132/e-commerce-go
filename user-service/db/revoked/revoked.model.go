package revoked

import "time"

type Revoked struct {
	ID    string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Value string `gorm:"index" json:"-"`

	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
