package admincontrollers

import (
	"time"

	basecontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
)

// IndexController 文档页面
type IndexController struct {
	basecontrollers.AdaptAdminController
}

// Index 文档主页
func (c *IndexController) Index() {
	c.AdaptAdminController.Data["pageTitle"] = "首页"
	c.AdaptAdminController.Data["ts"] = time.Now()
	c.AdaptAdminController.TplName = "index.html"
}
