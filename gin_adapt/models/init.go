package ginmodels

import (
	"fmt"
	"time"

	"github.com/lockeysama/go-easy-admin/geadmin/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db map[string]*gorm.DB

func initDB(name, dsn string) error {
	// dsn := "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
	if _db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	} else {
		if sqlDB, err := _db.DB(); err != nil {
			return err
		} else {
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(30)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}
		db[name] = _db
		return nil
	}
}

func DB(options ...string) *gorm.DB {
	if db == nil {
		db = make(map[string]*gorm.DB)
	}
	var (
		name string
		dsn  string
	)
	if len(options) == 0 {
		name = "default"
		dsn = utils.GetenvFromConfig("db.dsn", "").(string)
		if dsn == "" {
			panic("congig \"db.dsn\" can not be empty")
		}
		if err := initDB(name, dsn); err != nil {
			fmt.Println(err.Error())
			return nil
		}
	} else if len(options) == 1 {
		name = options[0]
	} else if len(options) == 2 {
		name, dsn = options[0], options[1]
		if err := initDB(name, dsn); err != nil {
			fmt.Println(err.Error())
			return nil
		}
	} else {
		panic("params exceptions")
	}
	if _db, ok := db[name]; ok {
		return _db
	}
	return nil
}
