package mysql

import (
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db           *sql.DB
	serverConfig map[string]interface{}
)

//var user string
//var password string
//var host string
//func InitMysql(u string, pw string, h string) {
//	if u == "" {
//		user = "root"
//	} else {
//		user = u
//	}
//
//	if pw == "" {
//		password = "Thunder#123"
//	} else {
//		password = pw
//	}
//
//	if h == "" {
//		host = "127.0.0.1"
//	} else {
//		host = h
//	}
//}

func LoadServerConfig() error {
	var opt ini.LoadOptions
	opt.IgnoreInlineComment = true
	cfg, err := ini.LoadSources(opt, "/opt/thunder/thunder.ini")
	if err != nil {
		fmt.Println(err)
		return err
	}
	serverConfig = map[string]interface{}{
		"host":     cfg.Section("MainServer").Key("DataBaseServerIp").String(),
		"passwd":   cfg.Section("MainServer").Key("Password").String(),
		"port":     "3306",
		"username": cfg.Section("MainServer").Key("UserName").String(),
		"dbname":   "karaok",
	}

	fmt.Println(serverConfig["passwd"])
	return nil
}

func GetDbInstance() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	//sqlCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	serverConfig["username"], serverConfig["passwd"], serverConfig["host"], serverConfig["port"], serverConfig["dbname"])

	sqlCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		serverConfig["username"], serverConfig["passwd"], serverConfig["host"], serverConfig["port"], serverConfig["dbname"])

	fmt.Println(sqlCon)

	db, err := sql.Open("mysql", sqlCon)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(15)
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
