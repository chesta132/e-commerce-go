package categorylib

import (
	"product-service/internal/model"
)

func FilterToCreate(payload model.CreateCategoryPayload) *model.Category {
	return &model.Category{
		Name:        payload.Name,
		Description: payload.Description,
	}
}
