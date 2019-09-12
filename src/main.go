package main

import (
	"http"
	"mysql"
	"os"
)

func main() {
	/*init mysql connect use gorm*/
	//mysql.InitMysql("root", "Thunder#123", "192.168.2.201")
	mysql.LoadServerConfig()

	/*create http web*/
	err := http.CreateHttpWeb()
	if err != nil {
		println(err)
		os.Exit(1)
	}
}
