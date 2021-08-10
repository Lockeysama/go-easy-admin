package geamodels

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Role))
}

// Role 角色
type Role struct {
	NormalModel
	Name         string   `orm:"unique" description:"名称" json:"name" display:"title=角色名"`
	Description  string   `description:"描述" json:"description" display:"title=描述;dbtype=Text"`
	Status       int      `description:"状态" json:"status" display:"title=状态"`
	CreatedAdmin *Admin   `orm:"rel(fk)" description:"创建者" json:"created_admin" display:"title=创建者;showfield=UserName"`
	UpdatedAdmin *Admin   `orm:"rel(fk)" description:"最后一次修改者" json:"updated_admin" display:"title=最后一次修改者;showfield=UserName"`
	Admins       []*Admin `orm:"-" description:"角色拥有者" json:"admins" display:"-"`
}

var (
	// RoleAdmin 超管角色
	RoleAdmin = "admin"
	// RoleUser 用户角色
	RoleUser = "user"
	// RoleAnonymous 匿名角色
	RoleAnonymous = "anonymous"
	// RolesID 角色 ID
	RolesID = map[string]int{
		RoleAdmin:     -1,
		RoleUser:      -1,
		RoleAnonymous: -1,
	}
)

// RegisterRoles 注册角色模型 - 初始化
func RegisterRoles() {
	o := orm.NewOrm()
	for key := range RolesID {
		name := GetRoleString(key)
		admin := Admin{UserName: "admin"}
		b := o.Read(&admin, "UserName")
		fmt.Println(b)
		r := Role{Name: name, CreatedAdmin: &admin, UpdatedAdmin: &admin}
		_, id, err := o.ReadOrCreate(&r, "Name")
		if err != nil {
			panic(err)
		}
		RolesID[key] = int(id)
	}
}

// GetRoleString 前缀 role_
func GetRoleString(s string) string {
	if strings.HasPrefix(s, "role_") {
		return s
	}
	return fmt.Sprintf("role_%s", s)
}

// AddRolesGroupPolicy 向Casbin添加角色继承策略规则
func AddRolesGroupPolicy() {
	// 普通管理员继承用户
	_, _ = Enforcer.AddGroupingPolicy(GetRoleString(RoleAdmin), GetRoleString(RoleUser))
	// 用户继承匿名者
	_, _ = Enforcer.AddGroupingPolicy(GetRoleString(RoleUser), GetRoleString(RoleAnonymous))
}
