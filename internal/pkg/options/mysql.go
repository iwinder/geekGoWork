package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlOption struct {
	Host                  string        `yaml:"host" mapstructure:"host"`
	Username              string        `yaml:"username" mapstructure:"username"`
	Password              string        `yaml:"password" mapstructure:"password"`
	Database              string        `yaml:"database" mapstructure:"database"`
	MaxIdleConnections    int           `yaml:"max-idle-connections" mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `yaml:"max-open-connections" mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max-connection-life-time" mapstructure:"max-connection-life-time"`
	LogLevel              int           `yaml:"log-level" mapstructure:"log-level"`
}

func New(opts *MysqlOption) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")
	// opts.LogLevel
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(opts.LogLevel)),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
