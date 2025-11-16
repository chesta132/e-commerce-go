package db

import (
	"product-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&model.Product{}, &model.ProductMeta{}, &model.Category{}, &model.StockMovement{}, &model.Image{}); err != nil {
		panic(err.Error())
	}

	db.Exec(`
        ALTER TABLE products 
        ADD COLUMN IF NOT EXISTS search_vector tsvector 
        GENERATED ALWAYS AS (
            to_tsvector('english', 
                coalesce(name,'') || ' ' || coalesce(description,'')
            )
        ) STORED
    `)

	db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_products_search 
        ON products USING GIN(search_vector)
    `)
}
