package admincontrollers

import (
	basecontrollers "TDCS/controllers/base"
)

// AdminBaseController 管理面板基础类
type AdminBaseController struct {
	basecontrollers.ManageBaseController
}

// Prefix 前缀
func (c *AdminBaseController) Prefix() string {
	return "/admin"
}

// PrefixAlias 前缀别名
func (c *AdminBaseController) PrefixAlias() string {
	return "系统管理"
}

// PrefixIcon 管理界面一级侧栏图标
func (c *AdminBaseController) PrefixIcon() string {
	return "layui-icon-set"
}
