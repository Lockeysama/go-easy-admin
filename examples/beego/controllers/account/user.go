package accountcontrollers

import (
	accountmodels "github.com/lockeysama/go-easy-admin/examples/beego/models/account"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// AccountController
type AccountController struct {
	AccountBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *AccountController) DBModel() geamodels.Model {
	return &accountmodels.User{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *AccountController) AdminNameAlias() string {
	return "用户"
}

// AdminIcon 管理界面侧栏图标
func (c *AccountController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *AccountController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAdminBaseController.GEADataBaseQueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*accountmodels.User)
	return x
}
