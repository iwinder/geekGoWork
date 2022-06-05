package mysql

import (
	"fmt"
	wderrors "github.com/iwinder/geekGoWork/internal/pkg/errors"
	configs "github.com/iwinder/geekGoWork/internal/pkg/options"
	perrors "github.com/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return perrors.Wrapf(err, fmt.Sprintf("DB Close error: %v", err))
	}
	return db.Close()
}

var (
	mysqlDb *datastore
	once    sync.Once
)

func NewMySQLDB(opts *configs.MysqlOption) *datastore {
	data, err := getMySQLDBOr(opts)
	if err != nil {
		fmt.Errorf("%d,MySql DB Open failed, error: %w", wderrors.ErrDatabase, err)
	}
	return data
}

func getMySQLDBOr(opts *configs.MysqlOption) (*datastore, error) {
	if opts == nil && mysqlDb.db == nil {
		return &datastore{}, fmt.Errorf("%d,MySql DB Open failed", wderrors.ErrDatabase)
	}
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = configs.New(opts)
		mysqlDb = &datastore{
			db: dbIns,
		}
	})
	if mysqlDb.db == nil || err != nil {
		return &datastore{}, perrors.Wrapf(err, "failed to get mysql store fatory, mysqlDB: %+v, error: %w", mysqlDb, err)
	}

	return mysqlDb, nil
}
