package geamodels

func InitModels() {
	RegisterCasbin()
	RegisterRoles()
	AddRolesGroupPolicy()
	CreateSuperUser()
}
