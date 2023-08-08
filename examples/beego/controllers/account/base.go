package accountcontrollers

import (
	basecontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/base"
)

// AccountBaseController 用户面板基础类
type AccountBaseController struct {
	basecontrollers.AdaptAdminController
}

// Prefix 前缀
func (c *AccountBaseController) Prefix() string {
	return "/account"
}

// PrefixAlias 前缀别名
func (c *AccountBaseController) PrefixAlias() string {
	return "账户管理"
}

// PrefixIcon 管理界面一级侧栏图标
func (c *AccountBaseController) PrefixIcon() string {
	return "layui-icon-set"
}
