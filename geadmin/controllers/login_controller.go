package geacontrollers

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"

	cache "github.com/lockeysama/go-easy-admin/geadmin/utils/cache"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// LoginController
type LoginController struct {
	BaseController
}

// Login 登录
// TODO:XSRF过滤
func (c *LoginController) Login() {
	if c.User != nil && c.User.ID > 0 {
		c.redirect("/home")
		return
	}
	// beego.ReadFromRequest(&c.Controller)
	if c.Ctx().RequestMethod() == "POST" {
		username := strings.TrimSpace(c.Ctx().InputQuery("username"))
		password := strings.TrimSpace(c.Ctx().InputQuery("password"))

		if username != "" && password != "" {
			var user = new(geamodels.Admin)
			o := orm.NewOrm()
			query := o.QueryTable(user)
			filters := map[string]interface{}{"username": username}
			for key := range filters {
				query = query.Filter(key, filters[key])
			}
			query.One(user)

			fmt.Println(user)
			// flash := beego.NewFlash()
			errorMsg := ""
			hash := md5.New()
			hash.Write([]byte(password + geamodels.Salt))
			if user == nil || user.Password != fmt.Sprintf("%x", hash.Sum(nil)) {
				errorMsg = "帐号或密码错误"
			} else if !user.Status {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIP = c.GetClientIP()
				user.LastLogin = time.Now().Unix()
				orm.NewOrm().
					QueryTable(user).
					Update(
						orm.Params{
							"LastLogin": time.Now().Unix(),
							"LastIP":    c.GetClientIP(),
						},
					)

				cache.DefaultMemCache().Set(
					fmt.Sprintf("uid%d", user.ID),
					user,
					cache.DefaultMemCacheExpiration,
				)
				hash := md5.New()
				hash.Write([]byte(c.GetClientIP() + "|" + user.Password + geamodels.Salt))
				authkey := fmt.Sprintf("%x", hash.Sum(nil))
				c.SetCookie("auth", fmt.Sprintf("%d|%s", user.ID, authkey), 7*86400)

				c.redirect("/home")
			}
			fmt.Println(errorMsg)
			// flash.Error(errorMsg)
			// flash.Store(&c.Controller)
			c.redirect("/login")
		}
	}
	c.TplName = "login/login.html"
}

// Logout 登出
func (c *LoginController) Logout() {
	c.SetCookie("auth", "")
	c.redirect("/login")
}

// NoAuth 无权限
func (c *LoginController) NoAuth() {
	c.WriteString("没有权限")
}
