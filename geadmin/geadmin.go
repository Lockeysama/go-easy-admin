package geadmin

import geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

type EngineType int8

const (
	Gin EngineType = iota
	Beego
)

func Router(engineType EngineType) map[string][]interface{} {
	switch engineType {
	case Gin:
	case Beego:
		return map[string][]interface{}{
			"Router": []interface{}{
				[]interface{}{"/", &geacontrollers.APIDocController{}, "*:Index"},
				[]interface{}{"/login", &geacontrollers.LoginController{}, "*:Login"},
				[]interface{}{"/logout", &geacontrollers.LoginController{}, "*:Logout"},
				[]interface{}{"/no_auth", &geacontrollers.LoginController{}, "*:NoAuth"},
				[]interface{}{"/home", &geacontrollers.HomeController{}, "*:Index"},
				[]interface{}{"/home/start", &geacontrollers.HomeController{}, "*:Start"},
			},
			"AutoRouter": []interface{}{
				&geacontrollers.APIDocController{},
			},
			"NSAutoRouter": []interface{}{
				"/admin",
				[]interface{}{&geacontrollers.AdminController{}},
				[]interface{}{&geacontrollers.RoleController{}},
				[]interface{}{&geacontrollers.CasbinController{}},
			}
		}
	}
}
