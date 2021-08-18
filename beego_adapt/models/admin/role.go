package adminmodels

import (
	"github.com/beego/beego/v2/client/orm"

	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Role))
}

// Role 角色
type Role struct {
	basemodels.ModelBase
	basemodels.NormalModel
	Name         string   `orm:"unique" description:"名称" json:"name" display:"title=角色名"`
	Description  string   `description:"描述" json:"description" display:"title=描述;dbtype=Text"`
	Status       int      `description:"状态" json:"status" display:"title=状态"`
	CreatedAdmin *Admin   `orm:"rel(fk)" description:"创建者" json:"created_admin" display:"title=创建者;showfield=UserName"`
	UpdatedAdmin *Admin   `orm:"rel(fk)" description:"最后一次修改者" json:"updated_admin" display:"title=最后一次修改者;showfield=UserName"`
	Admins       []*Admin `orm:"-" description:"角色拥有者" json:"admins" display:"-"`
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
	err = orm.NewOrm().QueryTable(&Role{}).Filter("Name", "role_"+geamodels.DefaultGEAdminUsername).One(role)
	return role, err
}
func (adapter *RoleAdapter) QueryRoleWithID(IDs ...int64) ([]geamodels.GEARole, error) {
	roles := new([]*Role)
	_, err := orm.NewOrm().QueryTable(&Role{}).Filter("id__in", IDs).All(roles)

	geaRoles := new([]geamodels.GEARole)
	for _, role := range *roles {
		*geaRoles = append(*geaRoles, role)
	}
	return *geaRoles, err
}
func (adapter *RoleAdapter) ReadOrCreate(
	role geamodels.GEARole, field string) (isCreate bool, ID int64, err error,
) {
	return orm.NewOrm().ReadOrCreate(role, field)
}
