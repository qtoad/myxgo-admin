package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qtoad/xgo-admin/modules/db"
)

var (
	orm *gorm.DB
	err error
)

func Init(c db.Connection) {
	orm, err = gorm.Open("sqlite3", c.GetDB("default"))

	if err != nil {
		panic("initialize orm failed")
	}
}
