package repository

import (
	"context"
	"github.com/EugeneNail/actum/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, name string, email string, password string) (id int64, err error)
	Find(ctx context.Context, column string, value any) (entity.User, error)
}
