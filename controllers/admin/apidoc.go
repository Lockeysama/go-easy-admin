package admincontrollers

import (
	basecontrollers "TDCS/controllers/base"
	"time"
)

// APIDocController 文档页面
type APIDocController struct {
	basecontrollers.BaseController
}

// Index 文档主页
func (c *APIDocController) Index() {
	c.Data["pageTitle"] = "API文档"
	c.Data["ts"] = time.Now()
	c.TplName = "apidoc/index.html"
}
