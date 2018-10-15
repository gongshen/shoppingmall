package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"shopping/models"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
	IsExit := c.Input().Get("exit") == "true"
	if IsExit {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("pwd", "", -1, "/")
		c.Redirect("/login", 302)
		return
	}
}

func (c *LoginController) Post() {
	username := c.Input().Get("uname")
	password := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLogin") == "on"
	//将password进行md5编码
	md5PWD := md5.New()
	io.WriteString(md5PWD, password)
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(buffer, "%x", md5PWD.Sum(nil))
	newPWD := buffer.String()
	//将username进行Base64编码后存进Cookie
	newUsername:=base64.StdEncoding.EncodeToString([]byte(username))
	user := models.GetUserInfo(username)
	if user.Username == username && user.Password == newPWD {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		c.Ctx.SetCookie("uname", newUsername, maxAge, "/")
		c.Ctx.SetCookie("pwd", newPWD, maxAge, "/")
	}
	/*设置session，暂时不用
	userInfo := models.GetUserInfo(username)
	if userInfo.Password == newPWD {
		sess, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
		sess.Set("uid", userInfo.Id)
		sess.Set("uname", userInfo.Username)
		fmt.Println("登陆成功！")
	} else {
		fmt.Println("登陆失败！")
	}
	*/
	c.Ctx.Redirect(302, "/")
	return
}

//对Cookie中的username进行Base64解码
func checkAccount(ctx *context.Context) bool {
	newUsername := ctx.GetCookie("uname")
	if newUsername==""{
		return false
	}
	pwd := ctx.GetCookie("pwd")
	username,_:=base64.StdEncoding.DecodeString(newUsername)
	uname:=string(username)
	userInfo := models.GetUserInfo(uname)
	return userInfo.Username == uname && userInfo.Password == pwd
}
