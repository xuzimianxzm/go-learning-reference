package main

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	name, _ := web.AppConfig.String("name")
	user, _ := web.AppConfig.String("mysqluser")
	password, _ := web.AppConfig.String("mysqlpassword")
	fmt.Println(name, user, password)

	web.Run()
}
