package adminmodels

import (
	"fmt"

	ginmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
	basemodels "github.com/lockeysama/go-easy-admin/gin_adapt/models/base"
)

func init() {
	if err := ginmodels.DB().AutoMigrate(&Role{}); err != nil {
		fmt.Println(err.Error())
	}
}

// Role 角色
type Role struct {
	basemodels.ModelBase   `gorm:"embedded"`
	basemodels.NormalModel `gorm:"embedded"`
	Name                   string `gorm:"unique" comment:"名称" json:"name" gea:"title=角色名"`
	Description            string `comment:"描述" json:"description" gea:"title=描述;dbtype=Text"`
	Status                 int    `comment:"状态" json:"status" gea:"title=状态"`
	CreatedAdminID         int64
	CreatedAdmin           *Admin `comment:"创建者" json:"created_admin" gea:"title=创建者;showfield=UserName"`
	UpdatedAdminID         int64
	UpdatedAdmin           *Admin   `comment:"最后一次修改者" json:"updated_admin" gea:"title=最后一次修改者;showfield=UserName"`
	Admins                 []*Admin `gorm:"-" description:"角色拥有者" json:"admins" gea:"-"`
}

func (m Role) TableName() string {
	return "admin_role"
}

func (m *Role) GetID() int64 {
	return m.ID
}

func (m *Role) GetName() string {
	return m.Name
}

func (m *Role) GetDescription() string {
	return m.Description
}

func (m *Role) GetStatus() int {
	return m.Status
}

func (m *Role) GetCreatedAdmin() geamodels.GEAdmin {
	return m.CreatedAdmin
}

func (m *Role) GetUpdatedAdmin() geamodels.GEAdmin {
	return m.UpdatedAdmin
}

func (m *Role) GetAdmins() []geamodels.GEAdmin {
	admins := new([]geamodels.GEAdmin)
	if m.Admins == nil {
		return *admins
	}

	for _, admin := range m.Admins {
		*admins = append(*admins, admin)
	}
	return *admins
}

type RoleAdapter struct{}

func (adapter *RoleAdapter) NewGEARole(name string, creator geamodels.GEAdmin) geamodels.GEARole {
	return &Role{Name: name, CreatedAdmin: creator.(*Admin), UpdatedAdmin: creator.(*Admin)}
}

func (adapter *RoleAdapter) QueryAdminRole() (role geamodels.GEARole, err error) {
	role = new(Role)
	result := ginmodels.DB().Model(&Role{}).First(role, "Name", "role_"+geamodels.DefaultGEAdminUsername)
	err = result.Error
	return role, err
}

func (adapter *RoleAdapter) QueryRoleWithID(IDs ...int64) ([]geamodels.GEARole, error) {
	roles := new([]*Role)
	result := ginmodels.DB().Find(roles, IDs)
	if result.Error != nil {
		return nil, result.Error
	}

	geaRoles := new([]geamodels.GEARole)
	for _, role := range *roles {
		*geaRoles = append(*geaRoles, role)
	}
	return *geaRoles, result.Error
}

func (adapter *RoleAdapter) ReadOrCreate(
	role geamodels.GEARole, field string) (isCreate bool, ID int64, err error,
) {
	result := ginmodels.DB().Model(&Role{}).FirstOrCreate(role, role)
	ID = role.GetID()
	return result.Error == nil, ID, result.Error
}
