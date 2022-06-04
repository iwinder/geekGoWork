package po

import (
	metaV1 "github.com/iwinder/geekGoWork/internal/pkg/meta/v1"
	"time"
)

type UserPO struct {
	metaV1.ObjectMeta `json:"metadata,omitempty"`
	Username          string     `json:"username,omitempty" gorm:"unique;colum:username;type:varchar(255);not null"`
	Nickname          string     `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`
	Password          string     `json:"password,omitempty" gorm:"column:password" validate:"required"`
	CreatedAt         time.Time  `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedAt         time.Time  `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	DeletedAt         *time.Time `json:"-" gorm:"column:deleted_at,index;comment:删除时间"`
}

func (u *UserPO) TableName() string {
	return "wd_sys_user"
}
