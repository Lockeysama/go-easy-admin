package admincontrollers

import (
	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// AdminController
type AdminController struct {
	basecontrollers.AdaptAdminController
}

// Index 主页
func (c *AdminController) Index() {
	c.Controller.Data["pageTitle"] = "系统首页"
	c.Controller.TplName = "admin/main.html"
}

// Start 控制面板
func (c *AdminController) Start() {
	c.Controller.Data["pageTitle"] = "控制面板"
	c.GEAdminBaseController.Display()
}
