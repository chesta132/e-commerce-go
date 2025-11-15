package repo

import (
	"context"
	"user-service/db/verification"

	"gorm.io/gorm"
)

type Verification struct {
	db *gorm.DB
}

func NewVerif(db *gorm.DB) *Verification {
	return &Verification{db}
}

func (r *Verification) FindByValue(ctx context.Context, value string) (verification.Verification, error) {
	return gorm.G[verification.Verification](r.db).Where("value = ?", value).First(ctx)
}

func (r *Verification) CreateOne(ctx context.Context, value, types string) (verification.Verification, error) {
	rev := verification.Verification{Value: value, Type: types}
	return rev, gorm.G[verification.Verification](r.db).Create(ctx, &rev)
}
