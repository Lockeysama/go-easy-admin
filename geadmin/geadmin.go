package geadmin

import geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

func AutoRegistryRouter(controller geacontrollers.GEARolePolicy) geacontrollers.GEARolePolicy {
	geacontrollers.RegisterGEARolePolicy(controller)
	geacontrollers.RegisterSideTree(controller)
	return controller
}
