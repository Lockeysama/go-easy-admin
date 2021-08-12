package geadmin

import geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

type EngineType string

const (
	Gin   EngineType = "Gin"
	Beego EngineType = "Beego"
	Other EngineType = "Other"
)

var EngineName EngineType

func InitEngine(engineType EngineType, baseController geacontrollers.Controller) {
	EngineName = engineType
	geacontrollers.EngineBaseController = baseController
}

type RouterType int8

const (
	Router RouterType = iota
	AutoRouter
)

func Routers() map[RouterType][]interface{} {
	return map[RouterType][]interface{}{
		Router: {
			[]interface{}{"/", &geacontrollers.APIDocController{}, "*:Index"},
			[]interface{}{"/login", &geacontrollers.LoginController{}, "*:Login"},
			[]interface{}{"/logout", &geacontrollers.LoginController{}, "*:Logout"},
			[]interface{}{"/no_auth", &geacontrollers.LoginController{}, "*:NoAuth"},
			[]interface{}{"/home", &geacontrollers.HomeController{}, "*:Index"},
			[]interface{}{"/home/start", &geacontrollers.HomeController{}, "*:Start"},
		},
		AutoRouter: {
			[]interface{}{"/admin", []interface{}{
				AutoRegistryRouter(&geacontrollers.AdminController{}),
				AutoRegistryRouter(&geacontrollers.RoleController{}),
				AutoRegistryRouter(&geacontrollers.CasbinController{}),
			}},
		},
	}
}

func AutoRegistryRouter(controller geacontrollers.ControllerRolePolicy) geacontrollers.ControllerRolePolicy {
	geacontrollers.RegisterControllerRolePolicy(controller)
	geacontrollers.RegisterSideTree(controller)
	return controller
}
