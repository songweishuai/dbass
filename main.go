package main

import (
	"dbass/http"
	"dbass/mysql"
	"fmt"
	"os"
)

func main() {
	fmt.Println("my name is dbass")

	/*init mysql connect use gorm*/
	mysql.InitMysql("root", "Thunder#123", "192.168.2.201")

	/*create http web*/
	err := http.CreateHttpWeb()
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
