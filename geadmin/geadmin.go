package geadmin

import geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"

type EngineType string

const (
	Gin   EngineType = "Gin"
	Beego EngineType = "Beego"
	Other EngineType = "Other"
)

var EngineName EngineType

func InitEngine(engineType EngineType, baseController geacontrollers.GEAController) {
	EngineName = engineType
	geacontrollers.EngineBaseController = baseController
}

func AutoRegistryRouter(controller geacontrollers.ControllerRolePolicy) geacontrollers.ControllerRolePolicy {
	geacontrollers.RegisterControllerRolePolicy(controller)
	geacontrollers.RegisterSideTree(controller)
	return controller
}
