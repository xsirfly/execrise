package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"exercise/conf"
)

var db *gorm.DB

func Init() error {
	dbconf := conf.GetConf().Database
	execriseDB, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.DBName))
	if err != nil {
		return err
	}
	db = execriseDB
	return nil
}
