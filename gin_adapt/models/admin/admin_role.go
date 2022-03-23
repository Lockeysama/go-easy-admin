package adminmodels

import (
	"fmt"

	ginmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models"
	"gorm.io/gorm"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	if err := ginmodels.DB().AutoMigrate(&AdminRole{}); err != nil {
		fmt.Println(err.Error())
	}
}

// AdminRole
type AdminRole struct {
	gorm.Model
	AdminID int64
	RoleID  int64
}

func (t AdminRole) TableName() string {
	return "admin_admin_role"
}

func (m *AdminRole) LoadM2M() {

}

func (m *AdminRole) GetID() int64 {
	return int64(m.ID)
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
	ginmodels.DB().Find(adminRoles, adminID)

	roles := new([]geamodels.GEAdminRole)
	for _, role := range *adminRoles {
		*roles = append(*roles, role)
	}
	return *roles
}

func (adapter *AdminRoleAdapter) Create(adminRole geamodels.GEAdminRole) (ID int64, err error) {
	row := new(AdminRole)
	row.ID = uint(adminRole.GetID())
	row.AdminID = adminRole.GetGEAdminID()
	row.RoleID = adminRole.GetGEARoleID()
	result := ginmodels.DB().Create(row)
	ID = int64(row.ID)
	err = result.Error
	return
}
