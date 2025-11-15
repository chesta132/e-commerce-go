package verification

import "time"

type Verification struct {
	ID    string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Value string `gorm:"index" json:"-"`
	Type  string `json:"-"`

	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

const (
	RequestToBecomeAdmin = "REQUEST_TO_BECOME_ADMIN"
)
