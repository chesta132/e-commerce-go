package repo

import (
	"context"
	"user-service/db/revoked"

	"gorm.io/gorm"
)

type Revoked struct {
	db *gorm.DB
}

func NewRevoked(db *gorm.DB) *Revoked {
	return &Revoked{db}
}

func (r *Revoked) FindByValue(ctx context.Context, value string) (revoked.Revoked, error) {
	return gorm.G[revoked.Revoked](r.db).Where("value = ?", value).First(ctx)
}

func (r *Revoked) CreateByValue(ctx context.Context, value string) (revoked.Revoked, error) {
	rev := revoked.Revoked{Value: value}
	return rev, gorm.G[revoked.Revoked](r.db).Create(ctx, &rev)
}
