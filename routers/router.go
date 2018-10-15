package routers

import (
	"github.com/astaxie/beego"
	"shopping/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/register",&controllers.RegisterControllers{})
	beego.Router("/showcart", &controllers.ShowController{}, "get:Showcart")
	beego.Router("/additem", &controllers.ShowController{}, "post:Additem")
	beego.Router("/addmin", &controllers.ShowController{}, "post:Addmin")
	beego.Router("/deleteitem", &controllers.ShowController{}, "post:Deleteitem")
	beego.Router("/order", &controllers.OrderController{})
	beego.AutoRouter(&controllers.OrderController{})
	beego.Router("/product",&controllers.ProductController{})
	beego.AutoRouter(&controllers.ProductController{})
}
