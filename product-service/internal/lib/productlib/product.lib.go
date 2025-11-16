package productlib

import "product-service/internal/model"

func FilterToCreate(payload *model.CreateProductPayload, adminId string, categories []model.Category) *model.Product {
	return &model.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Currency:    payload.Currency,
		AdminId:     adminId,
		Categories:  categories,
	}
}
