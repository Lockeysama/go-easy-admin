package applicationcontrollers

import (
	applicationmodels "github.com/lockeysama/go-easy-admin/examples/beego/models/application"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// ApplicationController
type ApplicationController struct {
	ApplicationBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *ApplicationController) DBModel() geamodels.Model {
	return &applicationmodels.Application{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *ApplicationController) AdminNameAlias() string {
	return "应用"
}

// AdminIcon 管理界面侧栏图标
func (c *ApplicationController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *ApplicationController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*applicationmodels.Application)
	return x
}
