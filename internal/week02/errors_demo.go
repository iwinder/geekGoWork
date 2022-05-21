package week02

import (
	"database/sql"
	"fmt"
	"github.com/iwinder/geekGoWork/internal/week02/service"
	"github.com/iwinder/geekGoWork/internal/week02/store"
	"github.com/iwinder/geekGoWork/internal/week02/store/mysql"
	"github.com/iwinder/geekGoWork/internal/week02/store/test"
	"github.com/pkg/errors"
)

// FindUser 调用 MySQL 查询
func FindUser() {
	dbf, err := mysql.GetMySQLDb()
	if err != nil {
		fmt.Printf("error is %+v\n", err)
		// 对于接口 可以返回 http.StatusInternalServerError，或者 http.StatusOK+自定义错误码+异常描述
		return
	}
	store.SetClient(dbf)
	user, errs := service.NewService(store.GetClient()).Users().QueryOne()
	if errs != nil {
		if errors.Is(errs, sql.ErrNoRows) {
			fmt.Println("The user info not found.")
			// 对于接口 可以返回 http.StatusNotFound，或者 http.StatusOK+自定义错误码+异常描述
			return
		} else {
			fmt.Printf("error is %+v\n", err)
			// 对于接口 可以返回 http.StatusInternalServerError，或者 http.StatusOK+自定义错误码+异常描述
			return
		}
	}
	// 对于接口 可以返回 http.StatusOK+自定义成功码+ user信息
	fmt.Printf("The user info %v", user)
	return

}

// FindUserBYTest 模拟查询用户信息
func FindUserBYTest() {
	dbf := test.GetTestDb()
	store.SetClient(dbf)
	user, errs := service.NewService(store.GetClient()).Users().QueryOne()
	if errs != nil {
		if errors.Is(errs, sql.ErrNoRows) {
			fmt.Println("The user info not found.")
			// 对于接口 可以返回 http.StatusNotFound，或者 http.StatusOK+自定义错误码+异常描述
			return
		} else {
			fmt.Printf("error is %+v\n", errs)
			// 对于接口 可以返回 http.StatusInternalServerError，或者 http.StatusOK+自定义错误码+异常描述
			return
		}
	}
	// 对于接口 可以返回 http.StatusOK+自定义成功码+ user信息
	fmt.Printf("The user info %v", user)
	return
}
