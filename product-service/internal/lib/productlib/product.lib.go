package productlib

import "product-service/internal/model"

func FilterToCreate(payload *model.CreateProductPayload, productMeta model.ProductMeta, adminId string, categories []model.Category) *model.Product {
	return &model.Product{
		Name:        payload.Name,
		Description: payload.Description,
		MetaId:      productMeta.ID,
		AdminId:     adminId,
		Categories:  categories,
	}
}

func FilterMetaToCreate(payload *model.CreateProductPayload) *model.ProductMeta {
	return &model.ProductMeta{
		Price:    payload.Price,
		Currency: payload.Currency,
	}
}
