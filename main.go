package main

import (
	db "github.com/debug-code/goDB/mysql"

	//"debug-code.com/models"
	_ "geekerblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func init() {
	mysqlParam := beego.AppConfig.String("mysqlParam")
	err := orm.RegisterDataBase("default",
		"mysql", mysqlParam, 30)
	if err != nil {
		beego.Error(err)
	}

	_, err = db.Regist("mysql",
		"alex:123qwe=-0@tcp(101.200.148.236:13306)/geekerblog?charset=utf8mb4",
		20, 50, time.Second*1800)
	if err != nil {
		beego.Error("mysql init", err)
	}
}

func main() {
	beego.Info(beego.BConfig.AppName, "V0.1")

	beego.Run()

}
