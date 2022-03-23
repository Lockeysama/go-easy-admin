package adminmodels

import (
	"fmt"

	ginmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models"
	"gorm.io/gorm"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	if err := ginmodels.DB().AutoMigrate(&Admin{}); err != nil {
		fmt.Println(err.Error())
	}
}

// Admin 管理员
type Admin struct {
	gorm.Model
	UserName  string  `gorm:"size:32" gea:"title=用户名"`
	Password  string  `gorm:"size:512" gea:"-"`
	RealName  string  `gorm:"size:32" gea:"title=真实姓名"`
	Phone     string  `gorm:"size:16" gea:"title=电话"`
	Email     string  `gorm:"size:32" gea:"title=电邮"`
	Avatar    string  `gorm:"size:256" gea:"title=头像;dbtype=File;required=false;meta=admin/avatar/"`
	Status    bool    `gea:"title=状态"`
	LastLogin int64   `gea:"title=最后登录时间;dbtype=Datetime"`
	LastIP    string  `gorm:"size:32" gea:"title=最后登录 IP"`
	Roles     []*Role `gorm:"-" gea:"title=拥有角色;showfield=Name"`
}

func (t Admin) TableName() string {
	return "admin_admin"
}

func (m *Admin) LoadM2M() {

}

func (m *Admin) GetID() int64 {
	return int64(m.ID)
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
	ginmodels.DB().Where("user_name", geamodels.DefaultGEAdminUsername).First(admin)
	return admin
}

func (adapter *AdminAdapter) QueryWithID(ID int64) geamodels.GEAdmin {
	admin := new(Admin)
	ginmodels.DB().First(admin, ID)
	return admin
}

func (adapter *AdminAdapter) ReadOrCreate(
	admin geamodels.GEAdmin, field string,
) (isCreate bool, ID int64, err error) {
	_admin := new(Admin)
	_admin.ID = uint(admin.GetID())
	_admin.UserName = admin.GetUserName()
	_admin.Password = admin.GetPassword()
	_admin.RealName = admin.GetRealName()
	_admin.Avatar = admin.GetAvatar()
	result := ginmodels.DB().FirstOrCreate(_admin, *_admin)
	ID = _admin.GetID()
	err = result.Error
	return
}
