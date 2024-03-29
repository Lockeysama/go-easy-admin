package admincontrollers

import (
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
	adminmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models/admin"
)

// AdministratorController
type AdministratorController struct {
	AdminBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *AdministratorController) DBModel() geamodels.Model {
	return &adminmodels.Admin{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *AdministratorController) AdminNameAlias() string {
	return "管理员"
}

// AdminIcon 管理界面侧栏图标
func (c *AdministratorController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *AdministratorController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.AdaptAdminController.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*adminmodels.Admin)
	return x
}
