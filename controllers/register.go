package controllers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	"io"
	"shopping/models"
)

var globalSessions *session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     36000,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

type RegisterControllers struct {
	beego.Controller
}

func (c *RegisterControllers) Get() {
	c.TplName = "register.html"
}

func (c *RegisterControllers) Post() {
	username := c.Input().Get("uname")
	password := c.Input().Get("pwd")
	telephone := c.Input().Get("tel")
	md5PWD := md5.New()
	io.WriteString(md5PWD, password)
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(buffer, "%x", md5PWD.Sum(nil))
	newPWD := buffer.String()
	userInfo := models.GetUserInfo(username)
	if userInfo.Username == "" {
		var user models.User
		_, err := user.Add(username, newPWD, telephone)
		if err != nil {
			fmt.Println("注册失败！")
		}
		sessn, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
		sessn.Set("uid", userInfo.Id)
		sessn.Set("uname", userInfo.Username)
		fmt.Println("注册成功！")
	} else {
		fmt.Println("用户已存在！")
	}
	c.Ctx.Redirect(302, "/")
	return
}
