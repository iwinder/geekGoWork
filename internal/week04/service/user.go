package service

import (
	"github.com/iwinder/geekGoWork/internal/week04/biz"
)

type UserService struct {
	biz *biz.UserBiz
}

func NewUserService(biz *biz.UserBiz) *UserService {
	return &UserService{
		biz: biz,
	}
}

//func (u *UserService) CreateUser(ctx context.Context, request *po.UserPO) {
//	u.biz.Create(ctx,request.)
//}
//
//func (u *UserService) GetUserInfo(ctx context.Context, request *po.UserPO) {
//	u.biz.Create(ctx,request.)
//}
