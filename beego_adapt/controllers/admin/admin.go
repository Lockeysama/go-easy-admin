package admincontrollers

import (
	"fmt"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// AdminController
type AdminController struct {
	AdminBaseController
}

// DBModel 返回控制器对应的数据库模型
func (c *AdminController) DBModel() geamodels.Model {
	return &geamodels.Admin{}
}

// AdminNameAlias 设置控制器侧栏别名
func (c *AdminController) AdminNameAlias() string {
	return "管理员"
}

// AdminIcon 管理界面侧栏图标
func (c *AdminController) AdminIcon() string {
	return "layui-icon-user"
}

func (c *AdminController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	l := c.GEAManageBaseController.QueryList(model, page, limit, filters, order, loadRel)
	x := l.(*[]*geamodels.Admin)
	fmt.Println(len(*x))
	return x
}
