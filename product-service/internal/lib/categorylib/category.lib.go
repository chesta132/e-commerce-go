package categorylib

import (
	"product-service/internal/model"
	"time"
)

func FilterToCreate(category *model.Category) *model.Category {
	if category == nil {
		return &model.Category{}
	}
	category.ID = ""
	category.CreatedAt = time.Time{}
	category.UpdatedAt = time.Time{}
	return category
}
