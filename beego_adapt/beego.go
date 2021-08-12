package beegoadapt

import (
	beego "github.com/beego/beego/v2/server/web"

	admincontrollers "github.com/lockeysama/go-easy-admin/beego_adapt/controllers/admin"
	"github.com/lockeysama/go-easy-admin/geadmin"
)

type RouterType int8

const (
	Router RouterType = iota
	AutoRouter
)

func Routers() map[RouterType][]interface{} {
	return map[RouterType][]interface{}{
		Router: {
			[]interface{}{"/", &admincontrollers.APIDocController{}, "*:Index"},
			[]interface{}{"/login", &admincontrollers.LoginController{}, "*:Login"},
			[]interface{}{"/logout", &admincontrollers.LoginController{}, "*:Logout"},
			[]interface{}{"/no_auth", &admincontrollers.LoginController{}, "*:NoAuth"},
			[]interface{}{"/home", &admincontrollers.HomeController{}, "*:Index"},
			[]interface{}{"/home/start", &admincontrollers.HomeController{}, "*:Start"},
		},
		AutoRouter: {
			[]interface{}{"/admin", []interface{}{
				geadmin.AutoRegistryRouter(&admincontrollers.AdminController{}),
				geadmin.AutoRegistryRouter(&admincontrollers.RoleController{}),
				geadmin.AutoRegistryRouter(&admincontrollers.CasbinController{}),
			}},
		},
	}
}

// Beego
func InitEngine() {
	// geadmin.InitEngine(geadmin.Beego, &AdaptController{})

	GEARouters := Routers()

	for routerType, routers := range GEARouters {
		for _, _router := range routers {
			router := _router.([]interface{})
			switch routerType {
			case Router:
				if len(router) != 3 {
					panic("router exception")
				}
				// controller := router[1].(geacontrollers.Controller)
				// controller.SetEngine(&AdaptController{})
				// controller := router[1].(beego.ControllerInterface)
				beego.Router(
					router[0].(string),
					router[1].(beego.ControllerInterface),
					router[2].(string),
				)
			case AutoRouter:
				group := router[0].(string)
				linkNamespaces := []beego.LinkNamespace{}
				for _, linkNamespace := range router[1].([]interface{}) {
					linkNamespaces = append(
						linkNamespaces,
						beego.NSAutoRouter(linkNamespace.(beego.ControllerInterface)),
					)
				}
				beego.AddNamespace(beego.NewNamespace(
					group,
					linkNamespaces...,
				))
			}
		}
	}
}
