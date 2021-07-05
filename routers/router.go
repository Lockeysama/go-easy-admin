// @APIVersion 1.0.0
// @Title Flowtime Test API
// @Description Flowtime Test API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	admincontrollers "TDCS/controllers/admin"
	basecontrollers "TDCS/controllers/base"
	datacontrollers "TDCS/controllers/data"
	hardwarecontrollers "TDCS/controllers/hardware"
	helpercontrollers "TDCS/controllers/helper"
	usercontrollers "TDCS/controllers/user"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &admincontrollers.APIDocController{}, "*:Index")
	beego.Router("/login", &admincontrollers.LoginController{}, "*:LoginIn")
	beego.Router("/login_out", &admincontrollers.LoginController{}, "*:LoginOut")
	beego.Router("/no_auth", &admincontrollers.LoginController{}, "*:NoAuth")

	beego.Router("/home", &admincontrollers.HomeController{}, "*:Index")
	beego.Router("/home/start", &admincontrollers.HomeController{}, "*:Start")
	beego.AutoRouter(&admincontrollers.APIDocController{})

	beego.AddNamespace(beego.NewNamespace("/admin",
		NSAutoRouter(&admincontrollers.AdminController{}),
		NSAutoRouter(&admincontrollers.RoleController{}),
		NSAutoRouter(&admincontrollers.CasbinController{}),
	))

	beego.AddNamespace(beego.NewNamespace("/user",
		NSAutoRouter(&usercontrollers.UserController{}),
		NSAutoRouter(&usercontrollers.AttributeController{}),
		NSAutoRouter(&usercontrollers.ConfigController{}),
		NSAutoRouter(&usercontrollers.DeviceUserController{}),
		NSAutoRouter(&usercontrollers.SocialController{}),
	))

	beego.AddNamespace(beego.NewNamespace("/data",
		NSAutoRouter(&datacontrollers.UsageRecordController{}),
	))

	beego.AddNamespace(beego.NewNamespace("/helper",
		NSAutoRouter(&helpercontrollers.IMHelperConfigController{}),
	))

	beego.AddNamespace(beego.NewNamespace("/hardware",
		NSAutoRouter(&hardwarecontrollers.HardwareConfigController{}),
	))

	basecontrollers.APIVersion = append(basecontrollers.APIVersion, "v1")

	userNS := beego.NewNamespace("/v1",
		beego.NSNamespace("/user/user", beego.NSInclude(&usercontrollers.UserController{})),
		beego.NSNamespace("/user/attribute", beego.NSInclude(&usercontrollers.AttributeController{})),
		beego.NSNamespace("/user/config", beego.NSInclude(&usercontrollers.ConfigController{})),
		beego.NSNamespace("/user/deviceuser", beego.NSInclude(&usercontrollers.DeviceUserController{})),
		beego.NSNamespace("/user/social", beego.NSInclude(&usercontrollers.SocialController{})),

		beego.NSNamespace("/data/log_file", beego.NSInclude(&datacontrollers.LogFileController{})),
		beego.NSNamespace("/data/usage_records", beego.NSInclude(&datacontrollers.UsageRecordController{})),

		beego.NSNamespace("/helper/im_helper_config", beego.NSInclude(&helpercontrollers.IMHelperConfigController{})),

		beego.NSNamespace("/hardware/config", beego.NSInclude(&hardwarecontrollers.HardwareConfigController{})),
	)
	beego.AddNamespace(userNS)
}

// AutoRouter 注册路由并注册 CasbinRule 和 SideTree
func AutoRouter(c basecontrollers.ControllerRolePolicy) *beego.HttpServer {
	basecontrollers.RegisterControllerRolePolicy(c)
	basecontrollers.RegisterSideTree(c)
	return beego.AutoRouter(c.(beego.ControllerInterface))
}

// NSAutoRouter 注册路由并注册 CasbinRule 和 SideTree
func NSAutoRouter(c basecontrollers.ControllerRolePolicy) beego.LinkNamespace {
	basecontrollers.RegisterControllerRolePolicy(c)
	basecontrollers.RegisterSideTree(c)
	return beego.NSAutoRouter(c.(beego.ControllerInterface))
}

// NSInclude 注册路由并注册 CasbinRule 和 SideTree
func NSInclude(c basecontrollers.ControllerRolePolicy) beego.LinkNamespace {
	basecontrollers.RegisterControllerRolePolicy(c)
	basecontrollers.RegisterSideTree(c)
	return beego.NSInclude(c.(beego.ControllerInterface))
}
