package adminmodels

import (
	"github.com/beego/beego/v2/client/orm"

	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(AdminRole))
}

// AdminRole
type AdminRole struct {
	basemodels.NormalModel
	AdminID int64 `orm:"column(admin_id)"`
	RoleID  int64 `orm:"column(role_id)"`
}

func (m *AdminRole) LoadM2M() {

}

func (m *AdminRole) GetID() int64 {
	return m.ID
}

func (m *AdminRole) GetGEAdminID() int64 {
	return m.AdminID
}

func (m *AdminRole) GetGEARoleID() int64 {
	return m.RoleID
}

type AdminRoleAdapter struct{}

func (adapter *AdminRoleAdapter) NewGEAdminRole(adminID int64, roleID int64) geamodels.GEAdminRole {
	return &AdminRole{AdminID: adminID, RoleID: roleID}
}

func (adapter *AdminRoleAdapter) QueryWithID(adminID int64) []geamodels.GEAdminRole {
	adminRoles := new([]*AdminRole)
	if _, err := orm.NewOrm().
		QueryTable(&AdminRole{}).
		Filter("AdminID", adminID).
		All(adminRoles); err != nil {
		return nil
	}

	roles := new([]geamodels.GEAdminRole)
	for _, role := range *adminRoles {
		*roles = append(*roles, role)
	}
	return *roles
}

func (adapter *AdminRoleAdapter) Create(adminRole geamodels.GEAdminRole) (ID int64, err error) {
	_, ID, err = orm.NewOrm().ReadOrCreate(adminRole, "AdminID", "RoleID")
	return
}
