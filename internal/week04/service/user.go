package service

import (
	"context"
	api "github.com/iwinder/geekGoWork/api/week04"
	"github.com/iwinder/geekGoWork/internal/week04/biz"
)

type UserService struct {
	api.UnimplementedUserServiceServer
	biz *biz.UserBiz
}

func NewUserService(biz *biz.UserBiz) *UserService {
	return &UserService{
		biz: biz,
	}
}

func (u *UserService) CreateUser(ctx context.Context, request *api.UserInfo) (*api.UserInfoReply, error) {
	err := u.biz.Create(ctx, &biz.UserDO{Username: request.Username,
		Nickname: request.Nickname,
		Password: request.Password,
	})
	return nil, err
}

func (u *UserService) GetUserInfo(ctx context.Context, request *api.UserInfo) (*api.UserInfoReply, error) {
	userDo, err := u.biz.GetOneById(ctx, request.Uid)
	if err != nil {

	}

	return &api.UserInfoReply{User: &api.UserRep{
		Username: userDo.Username,
		Nickname: userDo.Nickname,
		Uid:      userDo.ObjectMeta.ID,
	}}, nil

}
