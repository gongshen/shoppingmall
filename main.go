package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"shopping/models"
	_ "shopping/routers"
)

func init() {
	models.RegisterDB()
	orm.Debug = true
	orm.RunSyncdb("default", false, true)

}
func main() {
	os.Mkdir("attachment",os.ModePerm)
	beego.SetStaticPath("/attachment","attachment")
	beego.Run()
}
