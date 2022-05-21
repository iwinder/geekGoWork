package mysql

import (
	"database/sql"
	"github.com/iwinder/geekGoWork/internal/week02/store"
	perrors "github.com/pkg/errors"
	"sync"
	"time"
)

type datastore struct {
	db *sql.DB
}

// Users Mysql用于查询user信息的部分
func (ds *datastore) Users() store.UserStore {
	return newUserStore(ds)
}

func (ds *datastore) Close() {
	ds.db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLDb 获取M ySQL 数据库链接
func GetMySQLDb() (store.Factory, error) {
	var db *sql.DB
	var err error
	dns := "root:root@(127.0.0.1:3306)/testdatabases"
	once.Do(func() {
		db, err = sql.Open("mysql", dns)
		mysqlFactory = &datastore{db}
	})

	//defer db.Close() // 关闭数据库
	if mysqlFactory == nil || err != nil {
		return nil, perrors.Wrapf(err, "DB [%s] Open failed", dns)
	}
	err = db.Ping()
	if err != nil {
		return nil, perrors.Wrapf(err, "DB [%s] Ping failed", dns)
	}
	db.SetConnMaxLifetime(time.Second * 10) // 连接存活时间
	db.SetMaxOpenConns(200)                 // 最大连接数
	db.SetMaxIdleConns(10)                  // 最大空闲连接数

	return mysqlFactory, nil
}
