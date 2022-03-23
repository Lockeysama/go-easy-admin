package admincontrollers

import (
	"time"

	basecontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
)

// IndexController 文档页面
type IndexController struct {
	basecontrollers.AdaptController
}

// Index 文档主页
func (c *IndexController) Index() {
	c.AdaptController.Data["pageTitle"] = "首页"
	c.AdaptController.Data["ts"] = time.Now()
	c.AdaptController.TplName = "index.html"
}
