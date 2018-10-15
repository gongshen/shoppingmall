package controllers

import (
	"github.com/astaxie/beego"
	"shopping/helper"
	"shopping/models"
)

type OrderController struct {
	beego.Controller
}
func (c *OrderController) Prepare() {
	if !checkAccount(c.Ctx) {
		c.Ctx.Redirect(302, "/login")
		return
	}
}

func (c *OrderController) Get() {
	c.TplName = "order.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	uid := helper.Sti(c.Input().Get("u"))
	orders, err := models.GetAllOrders(uid)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Orders"] = orders
}

func (c *OrderController)Add(){
	uid := helper.Sti(c.GetString("uid"))
	ids := c.GetStrings("ids")
	totalp := helper.Stf64(c.GetString("totp"))
	balance:=models.GetBalanceById(uid)
	if balance<totalp{
		return
	}
	models.UpdateTheBalance(balance,totalp,uid)
	var order models.Order
	var remarks = "订单备注"
	var s = order.Addorder(ids, uid, remarks, totalp)
	resultjson := map[string]string{"status": s}
	c.Data["json"] = &resultjson
	c.ServeJSON()
}