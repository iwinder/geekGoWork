package service

import (
	"github.com/iwinder/geekGoWork/internal/week02/models"
	"github.com/iwinder/geekGoWork/internal/week02/store"
)

type UserSrv interface {
	QueryOne() (*models.User, error)
}

type userService struct {
	store store.Factory
}

// QueryOne 查询数据
func (u *userService) QueryOne() (*models.User, error) {
	return u.store.Users().QueryOne()
}

// 创建服务
func newUsers(srv *service) *userService {
	return &userService{
		store: srv.store,
	}
}
