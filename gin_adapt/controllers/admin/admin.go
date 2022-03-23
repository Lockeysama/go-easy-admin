package admincontrollers

import (
	basecontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
)

// AdminController
type AdminController struct {
	basecontrollers.AdaptController
}

// Index 主页
func (c *AdminController) Index() {
	c.AdaptController.Data["pageTitle"] = "系统首页"
	c.AdaptController.TplName = "admin/main.html"
}

// Start 控制面板
func (c *AdminController) Start() {
	c.AdaptController.Data["pageTitle"] = "控制面板"
	c.AdaptController.GEAdminBaseController.Display()
}
