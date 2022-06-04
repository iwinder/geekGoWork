package mysql

import (
	"fmt"
	"github.com/iwinder/geekGoWork/internal/week04/data"
	perrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

// UserRepo Mysql用于查询user信息的部分
func (ds *datastore) UserRepo() data.UserRepo {
	return newUserRepo(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return perrors.Wrapf(err, fmt.Sprintf("DB Close error: %v", err))
	}
	return db.Close()
}
