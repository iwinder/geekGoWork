package test

import (
	"github.com/iwinder/geekGoWork/internal/week02/store"
)

type datastore struct {
}

func (ds datastore) Users() store.UserStore {
	return newUserStore()
}

func GetTestDb() store.Factory {
	return &datastore{}
}
