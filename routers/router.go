// @APIVersion 1.0.0
// @Title Flowtime Test API
// @Description Flowtime Test API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &geacontrollers.APIDocController{}, "*:Index")
	beego.Router("/login", &geacontrollers.LoginController{}, "*:Login")
	beego.Router("/logout", &geacontrollers.LoginController{}, "*:Logout")
	beego.Router("/no_auth", &geacontrollers.LoginController{}, "*:NoAuth")

	beego.Router("/home", &geacontrollers.HomeController{}, "*:Index")
	beego.Router("/home/start", &geacontrollers.HomeController{}, "*:Start")
	beego.AutoRouter(&geacontrollers.APIDocController{})

	beego.AddNamespace(beego.NewNamespace("/admin",
		NSAutoRouter(&geacontrollers.AdminController{}),
		NSAutoRouter(&geacontrollers.RoleController{}),
		NSAutoRouter(&geacontrollers.CasbinController{}),
	))

	geacontrollers.APIVersion = append(geacontrollers.APIVersion, "v1")

	// userNS := beego.NewNamespace("/v1",
	// 	beego.NSNamespace("/user/user", beego.NSInclude(&usercontrollers.UserController{})),
	// )
	// beego.AddNamespace(userNS)
}

// AutoRouter 注册路由并注册 CasbinRule 和 SideTree
func AutoRouter(c geacontrollers.ControllerRolePolicy) *beego.HttpServer {
	geacontrollers.RegisterControllerRolePolicy(c)
	geacontrollers.RegisterSideTree(c)
	return beego.AutoRouter(c.(beego.ControllerInterface))
}

// NSAutoRouter 注册路由并注册 CasbinRule 和 SideTree
func NSAutoRouter(c geacontrollers.ControllerRolePolicy) beego.LinkNamespace {
	geacontrollers.RegisterControllerRolePolicy(c)
	geacontrollers.RegisterSideTree(c)
	return beego.NSAutoRouter(c.(beego.ControllerInterface))
}

// NSInclude 注册路由并注册 CasbinRule 和 SideTree
func NSInclude(c geacontrollers.ControllerRolePolicy) beego.LinkNamespace {
	geacontrollers.RegisterControllerRolePolicy(c)
	geacontrollers.RegisterSideTree(c)
	return beego.NSInclude(c.(beego.ControllerInterface))
}
