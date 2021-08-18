package geamodels

func InitModels() {
	RegisterCasbin()
	CreateAdministrator()
	RegisterRoles()
	AddRolesGroupPolicy()
	RegisterAdminRole()
}
