package geacontrollers

import (
	"path"
	"reflect"
	"strings"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// GEARolePolicy CasbinRule
type GEARolePolicy interface {
	DBModel() geamodels.Model
	Prefix() string

	PrefixAlias() string
	PrefixIcon() string
	AdminName() string
	AdminNameAlias() string
	AdminIcon() string
	AdminPathMethods() []string
}

// Prefix 前缀
func (c *GEAdminBaseController) Prefix() string {
	return "/"
}

// PrefixAlias 前缀别名
func (c *GEAdminBaseController) PrefixAlias() string {
	return ""
}

// PrefixIcon 管理界面一级侧栏图标（https://www.layui.com/doc/element/icon.html）
func (c *GEAdminBaseController) PrefixIcon() string {
	return ""
}

// AdminName 自定义控制器名称（默认使用控制器名称 Controller 的前面部分）
func (c *GEAdminBaseController) AdminName() string {
	return ""
}

// AdminNameAlias 自定义控制器名称（默认使用控制器名称 Controller 的前面部分）
func (c *GEAdminBaseController) AdminNameAlias() string {
	return ""
}

// AdminIcon 管理界面二级侧栏图标（https://www.layui.com/doc/element/icon.html）
func (c *GEAdminBaseController) AdminIcon() string {
	return ""
}

// AdminPathMethods 控制器 Admin 部分的请求函数（授权）
func (c *GEAdminBaseController) AdminPathMethods() []string {
	return []string{
		"list", "add", "update", "delete", "table", "edit",
		"ajaxadd", "ajaxupdate", "ajaxdelete",
		"ajaxupload", "ajaxdownload",
	}
}

// RegisterGEARolePolicy CasbinRule 和 SideNode）
func RegisterGEARolePolicy(r GEARolePolicy) {
	prefix := r.Prefix()
	adminPathMethods := r.AdminPathMethods()
	controllerName := r.AdminName()
	if controllerName == "" {
		reflectVal := reflect.ValueOf(r)
		ct := reflect.Indirect(reflectVal).Type()
		// TODO 丢出请使用 Controller 结尾或者重写 AdminName 方法
		controllerName = strings.TrimSuffix(ct.Name(), "Controller")
	}
	// 注册 Admin CasbinRule
	for _, _path := range adminPathMethods {
		pattern := path.Join(prefix, strings.ToLower(controllerName), _path)
		_, _ = geamodels.Enforcer.AddPolicy(geamodels.GetRoleString("admin"), pattern, "get")
	}
}
