package admincontrollers

import (
	adminmodels "TDCS/models/admin"
	basemodels "TDCS/models/base"
)

// AdminController
type AdminController struct {
	AdminBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *AdminController) DBModel() basemodels.Model {
	return &adminmodels.Admin{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *AdminController) AdminNameAlias() string {
	return "管理员"
}

// AdminIcon 管理界面侧栏图标
func (c *AdminController) AdminIcon() string {
	return "layui-icon-user"
}
