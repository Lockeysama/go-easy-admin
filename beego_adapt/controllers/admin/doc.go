package admincontrollers

import (
	"time"

	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// APIDocController 文档页面
type APIDocController struct {
	basecontrollers.AdaptController
}

// Index 文档主页
func (c *APIDocController) Index() {
	c.Controller.Data["pageTitle"] = "Docs"
	c.Controller.Data["ts"] = time.Now()
	c.Controller.TplName = "apidoc/index.html"
}
