package productlib

import (
	"fmt"
	"product-service/internal/model"
	"strings"

	"github.com/chesta132/e-commerce-go/shared/scrypto"
)

func FilterToCreate(payload *model.CreateProductPayload, productMeta model.ProductMeta, adminId string, categories []model.Category) *model.Product {
	if len(categories) < 1 {
		categories = append(categories, model.Category{Name: "Unknown", Description: "-"})
	}
	SKU := GenerateSKU(categories[0].Name, payload.Name)
	return &model.Product{
		Name:        payload.Name,
		Description: payload.Description,
		MetaId:      productMeta.ID,
		AdminId:     adminId,
		Categories:  categories,
		SKU:         SKU,
	}
}

func FilterMetaToCreate(payload *model.CreateProductPayload) *model.ProductMeta {
	return &model.ProductMeta{
		Price:    payload.Price,
		Currency: payload.Currency,
	}
}

func GenerateSKU(category, name string) string {
	cleanCat := strings.ToUpper(strings.ReplaceAll(category, " ", ""))
	prefixCat := cleanCat[:4]
	cleanName := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(cleanName) > 5 {
		cleanName = cleanName[:5]
	}

	randPart := strings.ToUpper(scrypto.RandomString(4))

	return fmt.Sprintf("%s-%s-%s", prefixCat, cleanName, randPart)
}
