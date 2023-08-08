package admincontrollers

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"

	adminmodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/admin"
)

// LoginController
type LoginController struct {
	basecontrollers.AdaptAdminController
}

// Login 登录
// TODO:XSRF过滤
func (c *LoginController) Login() {
	if c.User != nil && c.User.GetID() > 0 {
		c.Redirect("/admin", 302)
		return
	}
	beego.ReadFromRequest(&c.Controller)
	if c.RequestMethod() == "POST" {
		username := strings.TrimSpace(c.RequestQuery("username"))
		password := strings.TrimSpace(c.RequestQuery("password"))

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
			hash.Write([]byte(password + geamodels.Salt))
			if user == nil || user.Password != fmt.Sprintf("%x", hash.Sum(nil)) {
				errorMsg = "帐号或密码错误"
			} else if !user.Status {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIP = c.Controller.Ctx.Request.RemoteAddr
				user.LastLogin = time.Now().Unix()
				orm.NewOrm().
					QueryTable(user).
					Update(
						orm.Params{
							"LastLogin": time.Now().Unix(),
							"LastIP":    c.Controller.Ctx.Request.RemoteAddr,
						},
					)

				// cache.MemCache().Set(
				// 	fmt.Sprintf("uid%d", user.ID),
				// 	user,
				// )
				hash := md5.New()
				hash.Write([]byte(user.Password + geamodels.Salt))
				authkey := fmt.Sprintf("%x", hash.Sum(nil))
				c.SetCookie("auth", fmt.Sprintf("%d|%s", user.ID, authkey), 7*86400)

				c.Redirect("/admin", 302)
			}
			fmt.Println(errorMsg)
			flash.Error(errorMsg)
			flash.Store(&c.Controller)
			c.Redirect("/login", 302)
		}
	}
	c.Controller.TplName = "login/login.html"
}

// Logout 登出
func (c *LoginController) Logout() {
	c.SetCookie("auth", "")
	c.Redirect("/login", 302)
}

// NoAuth 无权限
func (c *LoginController) NoAuth() {
	c.Ctx.WriteString("没有权限")
}
