//go:build wireinject
// +build wireinject

package week04

import (
	"github.com/google/wire"
	configs "github.com/iwinder/geekGoWork/internal/pkg/options"
	"github.com/iwinder/geekGoWork/internal/week04/biz"
	"github.com/iwinder/geekGoWork/internal/week04/data/mysql"
	"github.com/iwinder/geekGoWork/internal/week04/service"
)

func InitUserServer(opts *configs.MysqlOption) *service.UserService {
	wire.Build(service.NewUserService, biz.NewUserBiz, mysql.NewUserRepo, mysql.NewMySQLDB)
	return &service.UserService{}
}
