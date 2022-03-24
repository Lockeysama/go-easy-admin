package adminmodels

import (
	"github.com/beego/beego/v2/client/orm"

	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Admin))
}

// Admin 管理员
type Admin struct {
	basemodels.NormalModel
	UserName  string  `orm:"unique" description:"用户名" gea:"title=用户名"`
	Password  string  `description:"密码" gea:"-"`
	RealName  string  `description:"真实姓名" gea:"title=真实姓名"`
	Phone     string  `description:"电话" gea:"title=电话"`
	Email     string  `description:"电邮" gea:"title=电邮"`
	Avatar    string  `gea:"title=头像;dbtype=File;required=false;meta=admin/avatar/"`
	Status    bool    `description:"状态" gea:"title=状态"`
	LastLogin int64   `orm:"auto_now" description:"最后登录时间" gea:"title=最后登录时间;dbtype=Datetime"`
	LastIP    string  `orm:"column(last_ip)" description:"最后登录 IP" gea:"title=最后登录 IP"`
	Roles     []*Role `orm:"reverse(many)" description:"拥有角色" gea:"title=拥有角色;showfield=Name"`
}

func (m *Admin) LoadM2M() {

}

func (m *Admin) GetID() int64 {
	return m.ID
}

func (m *Admin) GetUserName() string {
	return m.UserName
}

func (m *Admin) GetPassword() string {
	return m.Password
}

func (m *Admin) GetRealName() string {
	return m.RealName
}

func (m *Admin) GetAvatar() string {
	return m.Avatar
}

func (m *Admin) GetRoles() []geamodels.GEARole {
	roles := new([]geamodels.GEARole)

	if m.Roles == nil {
		return *roles
	}

	for _, role := range m.Roles {
		*roles = append(*roles, role)
	}
	return *roles
}

func (m *Admin) SetRoles(roles []geamodels.GEARole) {
	adminRoles := new([]*Role)
	for _, role := range roles {
		*adminRoles = append(*adminRoles, role.(*Role))
	}
	m.Roles = *adminRoles
}

type AdminAdapter struct{}

func (adapter *AdminAdapter) NewGEAdmin(username string, password string) geamodels.GEAdmin {
	return &Admin{UserName: username, Password: password, Status: true}
}

func (adapter *AdminAdapter) Administrator() geamodels.GEAdmin {
	admin := new(Admin)
	if err := orm.NewOrm().QueryTable(admin).
		Filter("UserName", geamodels.DefaultGEAdminUsername).
		One(admin); err != nil {
		return nil
	}
	return admin
}

func (adapter *AdminAdapter) QueryWithID(ID int64) geamodels.GEAdmin {
	admin := new(Admin)
	if err := orm.NewOrm().QueryTable(admin).Filter("ID", ID).One(admin); err != nil {
		return nil
	}
	return admin
}

func (adapter *AdminAdapter) ReadOrCreate(
	admin geamodels.GEAdmin, field string,
) (isCreate bool, ID int64, err error) {
	return orm.NewOrm().ReadOrCreate(admin, field)
}
