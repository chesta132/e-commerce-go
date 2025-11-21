package service

import (
	"context"
	"user-service/db/user"
	"user-service/internal/repo"

	"github.com/chesta132/e-commerce-go/shared/squery"
	"github.com/labstack/echo/v4"
)

type User struct {
	r *repo.User
}

type EchoUser struct {
	r   *repo.User
	c   echo.Context
	ctx context.Context
}

func NewUser(repo *repo.User) *User {
	return &User{repo}
}

func (r *User) AttachEcho(c echo.Context) *EchoUser {
	return &EchoUser{r: r.r, c: c, ctx: c.Request().Context()}
}

func (s *EchoUser) FindById(id string) (user.User, error) {
	return s.r.FindById(s.ctx, id)
}

func (s *EchoUser) UpdateById(id string, u user.User) error {
	return s.r.UpdateOne(s.ctx, []squery.Where{{Name: "id", Value: id}}, u)
}

func (s *EchoUser) FindByIdAndUpdate(id string, u user.User) (user.User, error) {
	err := s.UpdateById(id, u)
	if err != nil {
		return user.User{}, err
	}
	return s.r.FindById(s.ctx, id)
}
