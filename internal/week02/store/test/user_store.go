package test

import (
	"database/sql"
	"github.com/iwinder/geekGoWork/internal/week02/models"
	perrors "github.com/pkg/errors"
)

type userStore struct {
}

func newUserStore() *userStore {
	return &userStore{}
}

// 模拟数据库查询用户信息，返回  sql.ErrNoRows
func (u *userStore) sqlQuery() (*models.User, error) {
	return &models.User{}, sql.ErrNoRows
}

// QueryOne 查询单个用户信息
func (u *userStore) QueryOne() (*models.User, error) {
	user, err := u.sqlQuery()
	if err != nil {
		return user, perrors.Wrap(err, "DB found err")
	}

	return user, err
}
