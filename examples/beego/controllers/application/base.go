package applicationcontrollers

import (
	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// ApplicationBaseController 应用面板基础类
type ApplicationBaseController struct {
	basecontrollers.AdaptAdminController
}

// Prefix 前缀
func (c *ApplicationBaseController) Prefix() string {
	return "/application"
}

// PrefixAlias 前缀别名
func (c *ApplicationBaseController) PrefixAlias() string {
	return "应用管理"
}

// PrefixIcon 管理界面一级侧栏图标
func (c *ApplicationBaseController) PrefixIcon() string {
	return "layui-icon-set"
}
