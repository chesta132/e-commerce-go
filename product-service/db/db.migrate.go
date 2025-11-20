package db

import (
	"product-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&model.Product{}, &model.ProductMeta{}, &model.Category{}, &model.StockMovement{}, &model.Thumbnail{}); err != nil {
		panic(err.Error())
	}

	// enable trigram extension
	db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm;`)

	// drop dulu kl udh ada biar ga error
	db.Exec(`ALTER TABLE products DROP COLUMN IF EXISTS search_vector CASCADE;`)

	// column biasa dulu
	db.Exec(`ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector;`)

	// function trigger
	db.Exec(`
		CREATE OR REPLACE FUNCTION products_search_trigger() RETURNS trigger AS $$
		BEGIN
			new.search_vector := to_tsvector('english', coalesce(new.name,'') || ' ' || coalesce(new.description,''));
			RETURN new;
		END
		$$ LANGUAGE plpgsql;
	`)

	// trigger
	db.Exec(`DROP TRIGGER IF EXISTS tsvector_update ON products;`)
	db.Exec(`
		CREATE TRIGGER tsvector_update BEFORE INSERT OR UPDATE
		ON products FOR EACH ROW EXECUTE FUNCTION products_search_trigger();
	`)

	// index full text search
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_search_vector ON products USING GIN(search_vector);`)

	// index trigram (fuzzy search)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_name_trgm ON products USING gin(name gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_desc_trgm ON products USING gin(description gin_trgm_ops);`)
}
