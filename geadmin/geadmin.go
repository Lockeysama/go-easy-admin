package geadmin

import geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

func AutoRegistryRouter(controller geacontrollers.ControllerRolePolicy) geacontrollers.ControllerRolePolicy {
	geacontrollers.RegisterControllerRolePolicy(controller)
	geacontrollers.RegisterSideTree(controller)
	return controller
}
