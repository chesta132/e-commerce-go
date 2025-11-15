package repo

import (
	"context"
	"user-service/db/user"
	"user-service/internal/lib/query"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db}
}

func (r *User) FindById(ctx context.Context, id string) (user.User, error) {
	return gorm.G[user.User](r.db).Where("id = ?", id).First(ctx)
}

func (r *User) FindFirst(ctx context.Context, where []query.Where) (user.User, error) {
	q, v := query.BuildWhere(where)
	return gorm.G[user.User](r.db).Where(q, v...).First(ctx)
}

func (r *User) CreateOne(ctx context.Context, u *user.User) error {
	return gorm.G[user.User](r.db).Create(ctx, u)
}

func (r *User) UpdateOne(ctx context.Context, where []query.Where, u user.User) error {
	q, v := query.BuildWhere(where)
	_, err := gorm.G[user.User](r.db).Where(q, v...).Updates(ctx, u)
	return err
}
