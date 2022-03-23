package admincontrollers

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
	cache "github.com/lockeysama/go-easy-admin/geadmin/utils/cache"
	basecontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
	ginmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models"

	adminmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models/admin"
)

// LoginController
type LoginController struct {
	basecontrollers.AdaptController
}

// Login 登录
// TODO:XSRF过滤
func (c *LoginController) Login() {
	if c.User != nil && c.User.GetID() > 0 {
		c.Redirect("/admin", 302)
		return
	}
	// beego.ReadFromRequest(&c.Controller)

	if c.RequestMethod() == "POST" {
		username := strings.TrimSpace(c.RequestQuery("username"))
		password := strings.TrimSpace(c.RequestQuery("password"))

		if username != "" && password != "" {
			var user = new(adminmodels.Admin)
			filters := map[string]interface{}{"username": username}
			ginmodels.DB().Where(filters).First(user)

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
				user.LastIP = c.AdaptController.Ctx.Request.RemoteAddr
				user.LastLogin = time.Now().Unix()
				// orm.NewOrm().
				// 	QueryTable(user).
				// 	Update(
				// 		orm.Params{
				// 			"LastLogin": time.Now().Unix(),
				// 			"LastIP":    c.AdaptController.Ctx.Request.RemoteAddr,
				// 		},
				// 	)
				user.LastLogin = time.Now().Unix()
				user.LastIP = c.AdaptController.Ctx.Request.RemoteAddr
				ginmodels.DB().Updates(user)

				cache.MemCache().Set(
					fmt.Sprintf("uid%d", user.ID),
					user,
				)
				hash := md5.New()
				hash.Write([]byte(user.Password + geamodels.Salt))
				authkey := fmt.Sprintf("%x", hash.Sum(nil))
				c.SetCookie("auth", fmt.Sprintf("%d|%s", user.ID, authkey), 7*86400)

				c.Redirect("/admin", 302)
			}
			fmt.Println(errorMsg)
			// flash.Error(errorMsg)
			// flash.Store(&c.Controller)
			c.Redirect("/login", 302)
		}
	}
	c.AdaptController.TplName = "login/login.html"
}

// Logout 登出
func (c *LoginController) Logout() {
	c.SetCookie("auth", "")
	c.Redirect("/login", 302)
}

// NoAuth 无权限
func (c *LoginController) NoAuth() {
	c.Ctx.Writer.WriteString("没有权限")
}
