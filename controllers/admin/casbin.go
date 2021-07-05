package admincontrollers

import (
	adminmodels "TDCS/models/admin"
	basemodels "TDCS/models/base"
)

// CasbinController
type CasbinController struct {
	AdminBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *CasbinController) DBModel() basemodels.Model {
	return &adminmodels.CasbinRule{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *CasbinController) AdminNameAlias() string {
	return "CasbinRule"
}

// AdminIcon 管理界面侧栏图标
func (c *CasbinController) AdminIcon() string {
	return "layui-icon-key"
}
