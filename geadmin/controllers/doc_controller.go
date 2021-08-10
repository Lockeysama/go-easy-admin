package geacontrollers

import (
	"time"
)

// APIDocController 文档页面
type APIDocController struct {
	BaseController
}

// Index 文档主页
func (c *APIDocController) Index() {
	c.Data["pageTitle"] = "Docs"
	c.Data["ts"] = time.Now()
	c.TplName = "apidoc/index.html"
}
