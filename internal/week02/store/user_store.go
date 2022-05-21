package store

import "github.com/iwinder/geekGoWork/internal/week02/models"

// 定义查询用户信息的接口
type UserStore interface {
	// QueryOne 查询单个用户信息
	QueryOne() (*models.User, error)
}
