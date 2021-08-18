package geamodels

type GEAdminRoleAdapter interface {
	NewGEAdminRole(adminID int64, roleID int64) GEAdminRole

	QueryWithID(adminID int64) []GEAdminRole
	Create(adminRole GEAdminRole) (ID int64, err error)
}

var geadminRoleAdapter GEAdminRoleAdapter

func GetGEAdminRoleAdapter() GEAdminRoleAdapter {
	if geadminRoleAdapter == nil {
		panic("geadminRoleAdapter is nil")
	}
	return geadminRoleAdapter
}

func SetGEAdminRoleAdapter(adapter GEAdminRoleAdapter) {
	geadminRoleAdapter = adapter
}

type GEAdminRole interface {
	GetID() int64
	GetGEAdminID() int64
	GetGEARoleID() int64
}

func RegisterAdminRole() {
	if role, err := GetGEARoleAdapter().QueryAdminRole(); err != nil {
		panic(err.Error())
	} else {
		if role == nil {
			RegisterRoles()
			if role, err = GetGEARoleAdapter().QueryAdminRole(); err != nil || role == nil {
				panic("RegisterRoles failed")
			}
		}
		admin := GetGEAdminAdapter().Administrator()
		adminRole := GetGEAdminRoleAdapter().NewGEAdminRole(admin.GetID(), role.GetID())
		if _, err := GetGEAdminRoleAdapter().Create(adminRole); err != nil {
			panic(err.Error())
		}
	}
}
