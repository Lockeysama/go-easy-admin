package admincontrollers

import (
	"time"

	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// IndexController 文档页面
type IndexController struct {
	basecontrollers.AdaptAdminController
}

// Index 文档主页
func (c *IndexController) Index() {
	c.Controller.Data["pageTitle"] = "首页"
	c.Controller.Data["ts"] = time.Now()
	c.Controller.TplName = "index.html"
}
