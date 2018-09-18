package mysql

import (
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var user string
var password string
var host string

func InitMysql(u string, pw string, h string) {
	if u == "" {
		user = "root"
	} else {
		user = u
	}

	if pw == "" {
		password = "Thunder#123"
	} else {
		password = pw
	}

	if h == "" {
		host = "127.0.0.1"
	} else {
		host = h
	}
}

func GetDbInstance() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	sqlCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, "3306", "karaok")
	fmt.Println(sqlCon)

	db, err := sql.Open("mysql", sqlCon)
	if err != nil {
		return nil, err
	}

	//db.SetMaxOpenConns()
	//db.SetMaxIdleConns()
	//db.SetConnMaxLifetime()

	return db, nil
}

//func GetDbInstance() (*gorm.DB, error) {
//	if db != nil {
//		return db, nil
//	}
//
//	sql.Open()
//
//	sqlCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, "3306", "karaok")
//	fmt.Println(sqlCon)
//	//db, err := gorm.Open("mysql", "root:123456@/mysql")
//	db, err := gorm.Open("mysql", sqlCon)
//	if err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
