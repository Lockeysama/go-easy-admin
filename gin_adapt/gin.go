package ginadapt

import (
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
	admincontrollers "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/admin"
	gincontrolleradapt "github.com/lockeysama/go-easy-admin/gin_adapt/controllers/base"
	adminmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models/admin"
)

func InitGEAModelAdapt() {
	geamodels.SetGEACasbinAdapter(&adminmodels.CasbinRuleAdapter{})
	geamodels.SetGEARoleAdapter(&adminmodels.RoleAdapter{})
	geamodels.SetGEAdminAdapter(&adminmodels.AdminAdapter{})
	geamodels.SetGEAdminRoleAdapter(&adminmodels.AdminRoleAdapter{})
	geamodels.InitModels()
}

func InjectRouters() {
	gincontrolleradapt.RegisterRouter(&admincontrollers.IndexController{})
	// beego.Router("/", &admincontrollers.IndexController{}, "*:Index")

	// beego.Router("/login", &admincontrollers.LoginController{}, "*:Login")
	// beego.Router("/login_out", &admincontrollers.LoginController{}, "*:Logout")
	// beego.Router("/no_auth", &admincontrollers.LoginController{}, "*:NoAuth")

	gincontrolleradapt.RegisterRouter(&admincontrollers.AdministratorController{})
	gincontrolleradapt.RegisterRouter(&admincontrollers.RoleController{})
	gincontrolleradapt.RegisterRouter(&admincontrollers.CasbinController{})

	// beego.Router("/admin", &admincontrollers.AdminController{}, "*:Index")
	// beego.Router("/admin/start", &admincontrollers.AdminController{}, "*:Start")

	// beego.AddNamespace(beego.NewNamespace("/administrator",
	// 	AutoRegistryRouter(&admincontrollers.AdministratorController{}),
	// 	AutoRegistryRouter(&admincontrollers.RoleController{}),
	// 	AutoRegistryRouter(&admincontrollers.CasbinController{}),
	// ))
}
