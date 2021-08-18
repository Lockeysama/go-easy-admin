package admincontrollers

import (
	adminmodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/admin"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// RoleController
type RoleController struct {
	AdminBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *RoleController) DBModel() geamodels.Model {
	return &adminmodels.Role{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *RoleController) AdminNameAlias() string {
	return "角色"
}

// AdminIcon 管理界面侧栏图标
func (c *RoleController) AdminIcon() string {
	return "layui-icon-service"
}
