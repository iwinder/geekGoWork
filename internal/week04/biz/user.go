package biz

import (
	"context"
	metaV1 "github.com/iwinder/geekGoWork/internal/pkg/meta/v1"
	"github.com/iwinder/geekGoWork/internal/week04/data/mysql"
	"github.com/iwinder/geekGoWork/internal/week04/data/po"
	"time"
)

type UserBiz struct {
	repo *mysql.UserRepo
}

// NewUserBiz 新增
func NewUserBiz(repo *mysql.UserRepo) *UserBiz {
	return &UserBiz{
		repo: repo,
	}
}

func (u *UserBiz) Create(ctx context.Context, user *UserDO) error {
	userPO := &po.UserPO{
		Username: user.Username,
		Nickname: user.Nickname,
		Password: user.Password,
	}
	return u.repo.Create(ctx, userPO)
}

func (u *UserBiz) GetOneById(ctx context.Context, id uint64) (*UserDO, error) {
	userPO, err := u.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserDO{
		ObjectMeta: userPO.ObjectMeta,
		Username:   userPO.Username,
		Nickname:   userPO.Nickname,
		Password:   userPO.Password,
		CreatedAt:  userPO.CreatedAt,
	}, nil
}

type UserDO struct {
	metaV1.ObjectMeta `json:"metadata,omitempty"`
	Username          string    `json:"username,omitempty"`
	Nickname          string    `json:"nickname"`
	Password          string    `json:"password,omitempty"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
}
