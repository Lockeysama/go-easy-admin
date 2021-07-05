package admincontrollers

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	basecontrollers "TDCS/controllers/base"

	adminmodels "TDCS/models/admin"

	"TDCS/utils"

	cache "github.com/patrickmn/go-cache"
)

// LoginController
type LoginController struct {
	basecontrollers.BaseController
}

// LoginIn 登录
// TODO:XSRF过滤
func (c *LoginController) LoginIn() {
	if c.User != nil && c.User.ID > 0 {
		c.Redirect(beego.URLFor("HomeController.Index"))
	}
	beego.ReadFromRequest(&c.Controller)
	if c.Ctx.Request.Method == "POST" {

		username := strings.TrimSpace(c.GetString("username"))
		password := strings.TrimSpace(c.GetString("password"))

		if username != "" && password != "" {
			var user = new(adminmodels.Admin)
			o := orm.NewOrm()
			query := o.QueryTable(user)
			filters := map[string]interface{}{"username": username}
			for key := range filters {
				query = query.Filter(key, filters[key])
			}
			query.One(user)

			fmt.Println(user)
			flash := beego.NewFlash()
			errorMsg := ""
			hash := md5.New()
			hash.Write([]byte(password + adminmodels.Salt))
			if user == nil || user.Password != fmt.Sprintf("%x", hash.Sum(nil)) {
				errorMsg = "帐号或密码错误"
			} else if !user.Status {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIP = c.GetClientIP()
				user.LastLogin = time.Now().Unix()
				orm.NewOrm().QueryTable(user).Update(orm.Params{"LastLogin": time.Now().Unix(), "LastIP": c.GetClientIP()})

				utils.DefaultCache().Set("uid"+strconv.Itoa(user.ID), user, cache.DefaultExpiration)
				hash := md5.New()
				hash.Write([]byte(c.GetClientIP() + "|" + user.Password + adminmodels.Salt))
				authkey := fmt.Sprintf("%x", hash.Sum(nil))
				c.Ctx.SetCookie("auth", strconv.Itoa(user.ID)+"|"+authkey, 7*86400)

				c.Redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&c.Controller)
			c.Redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}
	c.TplName = "login/login.html"
}

// LoginOut 登出
func (c *LoginController) LoginOut() {
	c.Ctx.SetCookie("auth", "")
	c.Redirect(beego.URLFor("LoginController.LoginIn"))
}

// NoAuth 无权限
func (c *LoginController) NoAuth() {
	c.Ctx.WriteString("没有权限")
}
