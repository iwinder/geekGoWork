## 第二周作业
> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？


应该 Wrap 这个 error抛给上一层，因为此时在与第三方库(mysql、gorm等)进行协作，需要考虑保存堆栈信息。


###  [errors_demo_test.go](errors_demo_test.go):通过测试的方式模拟执行
- `TestFindUserByTest` 方法调用[errors_demo.go](errors_demo.go) 的`FindUserBYTest` 方法，不涉及数据库操作，直接模拟遇到 sql.ErrNoRows 的情况，最终打印异常 `The user info not found.`
  - 调用过程：  `TestFindUserByTest`->[errors_demo.go#FindUserBYTest](errors_demo.go)->[service/user_service.go#QueryOne](service/user_service.go)->[store/test/user_store.go#QueryOne](store/test/user_store.go)
```Go
// FindUserBYTest 模拟查询用户信息
func FindUserBYTest() {
	dbf := test.GetTestDb()
	store.SetClient(dbf)
	user, errs := service.NewService(store.GetClient()).Users().Query()
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
```
- `TestFindUserByMySql` 方法调用[errors_demo.go](errors_demo.go) 的`FindUser` 方法，涉及数据库操作，当不配置数据库信息时会因链接不上等原因报错，最终打印出堆栈信息
    - 调用过程：  `TestFindUserByTest`->[errors_demo.go#FindUser](errors_demo.go)->[service/user_service.go#QueryOne](service/user_service.go)->[store/mysql/user_store.go#QueryOne](store/mysql/user_store.go)
    - 由于目前关注点在异常信息处理方式，未提供数据库以及表结构等创建脚本
    - MySQL 数据库信息： `root:root@(127.0.0.1:3306)/testdatabases`
    - 查询语句 `select * from sys_user where id=1`
    - 涉及字段 [internal/week02/models/user.go](models/user.go)
## [service/user_service.go](service/user_service.go)
实际service，通过此处对象调用数据库查询部分，此处目前直接将数据库的数据抛到了上一层
```Go
func (u *userService) QueryOne() (*models.User, error) {
	return u.store.Users().QueryOne()
}
```
## [store](store) 数据库查询层
- [store/test/user_store.go](store/test/user_store.go) 模拟数据库查询返回 sql.ErrNoRows 的情况
```Go
// 模拟数据库查询用户信息，返回  sql.ErrNoRows
func (u *userStore) sqlQuery() (*models.User, error) {
return &models.User{}, sql.ErrNoRows
}

// QueryOne 查询单个用户信息，接收 sqlQuery 返回的信息，从而实现   sql.ErrNoRows 异常
func (u *userStore) QueryOne() (*models.User, error) {
	user, err := u.sqlQuery()
	if err != nil {
		return user, perrors.Wrap(err, "DB found err")
	}

	return user, err
}
```

- [store/mysql/user_store.go](store/mysql/user_store.go) 查询 MySQL 数据库的处理
```Go
// QueryOne  查询单个用户信息
func (u *userStore) QueryOne() (*models.User, error) {
	var user *models.User
	sql := "select * from sys_user where id=1"
	rows := u.db.QueryRow(sql)
	err := rows.Scan(&user)
	return user, perrors.Wrapf(err, fmt.Sprintf("sql: %s error: %v", sql, err))
}
```