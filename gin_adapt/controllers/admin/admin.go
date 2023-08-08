package admincontrollers

import (
	basecontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
)

// AdminController
type AdminController struct {
	basecontrollers.AdaptAdminController
}

// Index 主页
func (c *AdminController) Index() {
	c.AdaptAdminController.Data["pageTitle"] = "系统首页"
	c.AdaptAdminController.TplName = "admin/main.html"
}

// Start 控制面板
func (c *AdminController) Start() {
	c.AdaptAdminController.Data["pageTitle"] = "控制面板"
	c.AdaptAdminController.GEAdminBaseController.Display()
}
