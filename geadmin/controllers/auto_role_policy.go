package geacontrollers

import (
	"path"
	"reflect"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// ControllerRolePolicy CasbinRule
type ControllerRolePolicy interface {
	beego.ControllerInterface
	DBModel() geamodels.Model
	Prefix() string
	PrefixAlias() string
	PrefixIcon() string
	AdminName() string
	AdminNameAlias() string
	AdminIcon() string
	AdminPathMethods() []string
	RESTFulAPIPathMethods() map[string][]string
}

// Prefix 前缀
func (c *BaseController) Prefix() string {
	return "/"
}

// PrefixAlias 前缀别名
func (c *BaseController) PrefixAlias() string {
	return ""
}

// AdminName 自定义控制器名称（默认使用控制器名称 Controller 的前面部分）
func (c *BaseController) AdminName() string {
	return ""
}

// AdminNameAlias 自定义控制器名称（默认使用控制器名称 Controller 的前面部分）
func (c *BaseController) AdminNameAlias() string {
	return ""
}

// AdminPathMethods 控制器 Admin 部分的请求函数（授权）
func (c *BaseController) AdminPathMethods() []string {
	return []string{"list", "add", "update", "delete", "ajaxsave", "ajaxdel", "ajaxupload", "ajaxgetfile", "table", "edit"}
}

// RESTFulAPIPathMethods 控制器 RESTFul API 部分的请求函数（授权）
func (c *BaseController) RESTFulAPIPathMethods() map[string][]string {
	return map[string][]string{
		"/":      {"get", "post"},
		"/{id}/": {"get", "put", "patch", "delete"},
	}
}

// RegisterControllerRolePolicy CasbinRule 和 SideNode）
func RegisterControllerRolePolicy(r ControllerRolePolicy) {
	prefix := r.Prefix()
	adminPathMethods := r.AdminPathMethods()
	restAPIPathMethods := r.RESTFulAPIPathMethods()
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
	// 注册 RESTFul API CasbinRule
	for _path, methods := range restAPIPathMethods {
		for _, method := range methods {
			pattern := path.Join(prefix, strings.ToLower(controllerName), _path)
			_, _ = geamodels.Enforcer.AddPolicy(geamodels.GetRoleString("admin"), pattern, method)
		}
	}
}
