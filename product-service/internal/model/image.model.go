package model

type Image struct {
	ID        string `gorm:"primarykey;default:gen_random_uuid()" json:"id"`
	Image     []byte `gorm:"type:bytea" json:"image"`
	ProductId string `json:"productId"`
}
