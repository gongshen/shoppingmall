package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"shopping/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Prepare() {
	if !checkAccount(c.Ctx) {
		c.Ctx.Redirect(302, "/login")
		return
	}
}

func (c *MainController) Get() {
	c.TplName = "show.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	uname,_:=base64.StdEncoding.DecodeString(c.Ctx.GetCookie("uname"))
	userInfo:=models.GetUserInfo(string(uname))
	c.Data["Uid"]=userInfo.Id
	c.Data["Balance"]=userInfo.Balance
	products, err := models.GetAllProduct()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Products"] = products
}
