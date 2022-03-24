package applicationmodels

import (
	"github.com/beego/beego/v2/client/orm"

	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

func init() {
	orm.RegisterModelWithPrefix("application_", new(Application))

	geacontrollers.DefaultValueMakerRegister("AppKeyMaker", AppKeyMaker)
}

// User 用户
type Application struct {
	basemodels.NormalModel
	Name   string `orm:"size(32);description(名称)"`
	NameEN string `orm:"size(32);description(英文名称)"`
	Desc   string `orm:"size(256);description(描述)"`
	Status int8   `orm:"description(应用状态> 0: 已停止 | 1: 开发中 | 2: 已上线)"`
	AppKey string `orm:"size(36)" gea:"maker=AppKeyMaker"`
	Secret string `orm:"size(32)"`
}

func AppKeyMaker() interface{} {
	return "hello-app"
}
