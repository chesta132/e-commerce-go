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
	prefixCat := strings.ToUpper(category[:4])
	cleanName := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(cleanName) > 5 {
		cleanName = cleanName[:5]
	}

	randPart := strings.ToUpper(scrypto.RandomString(4))

	return fmt.Sprintf("%s-%s-%s", prefixCat, cleanName, randPart)
}

func AddRelationIdToFlat(product *model.Product) {
	for _, p := range product.Previews {
		product.PreviewIds = append(product.PreviewIds, p.ID)
	}
	for _, c := range product.Categories {
		product.CategoryIds = append(product.CategoryIds, c.ID)
	}
}

func MoveRelationIdToFlat(product *model.Product) {
	AddRelationIdToFlat(product)
	product.Previews = nil
	product.Categories = nil
}
