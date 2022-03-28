package applicationcontrollers

import (
	applicationmodels "github.com/lockeysama/go-easy-admin/examples/beego/models/application"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// PlatformController
type PlatformController struct {
	ApplicationBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *PlatformController) DBModel() geamodels.Model {
	return &applicationmodels.Platform{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *PlatformController) AdminNameAlias() string {
	return "Platform"
}

// AdminIcon 管理界面侧栏图标
func (c *PlatformController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *PlatformController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*applicationmodels.Platform)
	return x
}
