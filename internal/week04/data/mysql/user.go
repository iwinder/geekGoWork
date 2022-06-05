package mysql

import (
	"context"
	"errors"
	"fmt"
	wderrors "github.com/iwinder/geekGoWork/internal/pkg/errors"
	"github.com/iwinder/geekGoWork/internal/week04/data/po"
	"gorm.io/gorm"
	"regexp"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(ds *datastore) *UserRepo {
	return &UserRepo{
		db: ds.db,
	}
}

func (u *UserRepo) Create(ctx context.Context, user *po.UserPO) error {
	err := u.db.Create(&user).Error
	if err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			return fmt.Errorf("%d,data had crated ,error is [%s]:", wderrors.ErrAlreadyExist, err)
		}
		return fmt.Errorf("%d,data crate error,error is [%s]:", wderrors.ErrDatabase, err)
	}
	return err
}

func (u *UserRepo) GetOne(ctx context.Context, id uint64) (*po.UserPO, error) {
	user := &po.UserPO{}
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%d,data not found,error is [%s]:", wderrors.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%d,data query error,error is [%s]:", wderrors.ErrDatabase, err)
	}

	return user, nil
}
