package data

import (
	"context"
	"github.com/iwinder/geekGoWork/internal/week04/data/po"
)

type UserRepo interface {
	Create(ctx context.Context, user *po.UserPO) error
	GetOne(ctx context.Context, id uint64) (*po.UserPO, error)
}
