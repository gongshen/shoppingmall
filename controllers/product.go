package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"shopping/helper"
	"shopping/models"
)

type ProductController struct {
	beego.Controller
}

func (c *ProductController) Prepare() {
	if !checkAccount(c.Ctx) {
		c.Ctx.Redirect(302, "/login")
		return
	}
}

func (c *ProductController) Post() {
	name := c.Input().Get("name")
	price := helper.Stf64(c.Input().Get("price"))
	stock := helper.Sti(c.Input().Get("stock"))
	content := c.Input().Get("content")
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		attachment = fh.Filename
		beego.Info("Attachment:",attachment )
		err = c.SaveToFile("attachment", path.Join("attachment",attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	var p models.Product
	err = p.AddProduct(name, price, stock, content, attachment)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/", 302)
	return
}

func (c *ProductController) Add() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.TplName = "add_product.html"
}
