package admincontrollers

import (
	"fmt"

	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// HomeController
type HomeController struct {
	basecontrollers.AdaptController
}

// Index 主页
func (c *HomeController) Index() {
	c.Controller.Data["pageTitle"] = "系统首页"
	c.Controller.TplName = "public/main.html"
	fmt.Println("ss")
}

// Start 控制面板
func (c *HomeController) Start() {
	c.Controller.Data["pageTitle"] = "控制面板"
	c.GEAdminBaseController.Display()
}
