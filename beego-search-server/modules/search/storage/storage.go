package searchstorage

import "github.com/beego/beego/v2/client/orm"

type storage struct {
	db orm.Ormer
}

func NewStorage(db orm.Ormer) *storage {
	return &storage{db: db}
}
