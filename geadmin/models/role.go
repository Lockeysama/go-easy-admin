package geamodels

import (
	"fmt"
	"strings"
)

type GEARoleAdapter interface {
	NewGEARole(name string, creator GEAdmin) GEARole

	QueryAdminRole() (GEARole, error)
	QueryRoleWithID(...int64) ([]GEARole, error)
	ReadOrCreate(role GEARole, field string) (isCreate bool, ID int64, err error)
}

var gearoleAdapter GEARoleAdapter

func GetGEARoleAdapter() GEARoleAdapter {
	if gearoleAdapter == nil {
		panic("gearoleAdapter is nil")
	}
	return gearoleAdapter
}

func SetGEARoleAdapter(adapter GEARoleAdapter) {
	gearoleAdapter = adapter
}

// Role 角色
type GEARole interface {
	GetID() int64
	GetName() string
	GetDescription() string
	GetStatus() int
	GetCreatedAdmin() GEAdmin
	GetUpdatedAdmin() GEAdmin
	GetAdmins() []GEAdmin
}

type GEARoleType string

const (
	// RoleAdmin 超管角色
	RoleAdmin GEARoleType = "admin"
	// RoleUser 用户角色
	RoleUser GEARoleType = "user"
	// RoleAnonymous 匿名角色
	RoleAnonymous GEARoleType = "anonymous"
)

// RolesID 角色 ID
var RolesID = map[GEARoleType]int{
	RoleAdmin:     -1,
	RoleUser:      -1,
	RoleAnonymous: -1,
}

// RegisterRoles 注册角色模型 - 初始化
func RegisterRoles() {
	for key := range RolesID {
		name := GetRoleString(key)
		admin := GetGEAdminAdapter().NewGEAdmin(DefaultGEAdminUsername, "")
		if _, _, err := GetGEAdminAdapter().ReadOrCreate(admin, "UserName"); err != nil {
			panic("error")
		} else {
			r := GetGEARoleAdapter().NewGEARole(name, admin)
			_, id, err := GetGEARoleAdapter().ReadOrCreate(r, "Name")
			if err != nil {
				panic(err)
			}
			RolesID[key] = int(id)
		}
	}
}

// GetRoleString 前缀 role_
func GetRoleString(s GEARoleType) string {
	if strings.HasPrefix((string)(s), "role_") {
		return (string)(s)
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
