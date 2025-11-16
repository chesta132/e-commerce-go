package model

type CreateProductPayload struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Price       int         `json:"price" validate:"required"`
	Currency    string      `json:"currency" validate:"required"`
	Categories  *[]Category `json:"categories" validate:"required"`
	Image       *Image      `json:"image"`
}
