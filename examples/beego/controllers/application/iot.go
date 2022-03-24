package applicationcontrollers

import (
	applicationmodels "github.com/lockeysama/go-easy-admin/examples/beego/models/application"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// IoTController
type IoTController struct {
	ApplicationBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *IoTController) DBModel() geamodels.Model {
	return &applicationmodels.IoT{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *IoTController) AdminNameAlias() string {
	return "IoT"
}

// AdminIcon 管理界面侧栏图标
func (c *IoTController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *IoTController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*applicationmodels.IoT)
	return x
}
