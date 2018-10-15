package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"shopping/helper"
	"shopping/models"
)

type ShowController struct {
	beego.Controller
}

func (c *ShowController) Prepare() {
	if !checkAccount(c.Ctx) {
		c.Ctx.Redirect(302, "/login")
		return
	}
}

func (c *ShowController) Additem() {
	pid := helper.Sti(c.GetString("p"))
	num := helper.Sti(c.GetString("n"))
	uid := helper.Sti(c.GetString("u"))
	remarks:=c.GetString("remarks")
	var item models.Item
	var s = item.Additem(pid, uid, num,remarks)
	resultjson := map[string]string{"status": s}
	c.Data["json"] = &resultjson
	c.ServeJSON()
}

func (c *ShowController) Showcart() {
	uid := helper.Sti(c.GetString("u"))
	u := new(models.User)
	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", uid).One(u)
	c.Data["Balance"] = u.Balance
	c.Data["Uid"] = uid
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	var maps []orm.Params = models.Queryitems(uid)
	c.Data["list"] = maps
	c.TplName = "mycart.html"
}
func (c *ShowController) Addmin() {
	pid := helper.Sti(c.GetString("p"))
	num := helper.Sti(c.GetString("n"))
	var item models.Itemnum
	var s = item.Addmin(pid, num)
	resultjson := map[string]string{"status": s}
	c.Data["json"] = &resultjson
	c.ServeJSON()
}
func (c *ShowController) Deleteitem() {
	pid := helper.Sti(c.GetString("p"))
	var item models.Item
	var s = item.Deleteitem(pid)
	resultjson := map[string]string{"status": s}
	c.Data["json"] = &resultjson
	c.ServeJSON()
}


