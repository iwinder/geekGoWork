package mysql

import (
	"database/sql"
	"fmt"
	"github.com/iwinder/geekGoWork/internal/week02/models"
	perrors "github.com/pkg/errors"
)

type userStore struct {
	db *sql.DB
}

// QueryOne  查询单个用户信息
func (u *userStore) QueryOne() (*models.User, error) {
	var user *models.User
	sql := "select * from sys_user where id=1"
	rows := u.db.QueryRow(sql)
	err := rows.Scan(&user)
	return user, perrors.Wrapf(err, fmt.Sprintf("sql: %s error: %v", sql, err))
}

func newUserStore(ds *datastore) *userStore {
	return &userStore{
		db: ds.db,
	}
}
