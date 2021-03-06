package beegoadapt

import (
	beego "github.com/beego/beego/v2/server/web"

	admincontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/admin"
	adminmodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/admin"
	"github.com/lockeysama/go-easy-admin/geadmin"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

type RouterType int8

const (
	Router RouterType = iota
	AutoRouter
)

func InitGEAModelAdapt() {
	geamodels.SetGEACasbinAdapter(&adminmodels.CasbinRuleAdapter{})
	geamodels.SetGEARoleAdapter(&adminmodels.RoleAdapter{})
	geamodels.SetGEAdminAdapter(&adminmodels.AdminAdapter{})
	geamodels.SetGEAdminRoleAdapter(&adminmodels.AdminRoleAdapter{})
	geamodels.InitModels()
}

func InjectRouters() {
	beego.Router("/", &admincontrollers.IndexController{}, "*:Index")

	beego.Router("/login", &admincontrollers.LoginController{}, "*:Login")
	beego.Router("/login_out", &admincontrollers.LoginController{}, "*:Logout")
	beego.Router("/no_auth", &admincontrollers.LoginController{}, "*:NoAuth")

	beego.Router("/admin", &admincontrollers.AdminController{}, "*:Index")
	beego.Router("/admin/start", &admincontrollers.AdminController{}, "*:Start")

	beego.AddNamespace(beego.NewNamespace("/administrator",
		AutoRegistryRouter(&admincontrollers.AdministratorController{}),
		AutoRegistryRouter(&admincontrollers.RoleController{}),
		AutoRegistryRouter(&admincontrollers.CasbinController{}),
	))
}

func AutoRegistryRouter(controller geacontrollers.GEARolePolicy) beego.LinkNamespace {
	return beego.NSAutoRouter(
		geadmin.AutoRegistryRouter(controller).(beego.ControllerInterface),
	)
}
