package accountcontrollers

import (
	accountmodels "github.com/lockeysama/go-easy-admin/examples/beego/models/account"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// IAMController
type IAMController struct {
	AccountBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *IAMController) DBModel() geamodels.Model {
	return &accountmodels.IAM{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *IAMController) AdminNameAlias() string {
	return "IAM"
}

// AdminIcon 管理界面侧栏图标
func (c *IAMController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *IAMController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*accountmodels.IAM)
	return x
}
