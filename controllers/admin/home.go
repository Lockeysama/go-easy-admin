package admincontrollers

import basecontrollers "TDCS/controllers/base"

// HomeController
type HomeController struct {
	basecontrollers.BaseController
}

// Index 主页
func (c *HomeController) Index() {
	c.Data["pageTitle"] = "系统首页"
	c.TplName = "public/main.html"
}

// Start 控制面板
func (c *HomeController) Start() {
	c.Data["pageTitle"] = "控制面板"
	c.Display()
}
