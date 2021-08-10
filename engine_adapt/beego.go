package geaadapt

import (
	beego "github.com/beego/beego/v2/server/web"

	"github.com/lockeysama/go-easy-admin/geadmin"
)

// Beego
func InitEngine() {
	geadmin.InitEngine(geadmin.Beego, beego.Controller{})

	GEARouters := geadmin.Routers()

	for routerType, routers := range GEARouters {
		for _, _router := range routers {
			router := _router.([]interface{})
			switch routerType {
			case geadmin.Router:
				if len(router) != 3 {
					panic("router exception")
				}
				beego.Router(
					router[0].(string),
					router[1].(beego.ControllerInterface),
					router[2].(string),
				)
			case geadmin.AutoRouter:
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
