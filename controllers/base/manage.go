package basecontrollers

import (
	basemodels "TDCS/models/base"
	"bytes"

	"encoding/gob"
	"reflect"

	"github.com/beego/beego/v2/server/web/context"
)

// ManageBaseController 控制器管理基类
type ManageBaseController struct {
	BaseController
	Model     basemodels.Model
	PageTitle string
}

// Init 初始化
func (c *ManageBaseController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Instance = app.(ControllerRolePolicy)
	c.Model = c.Instance.DBModel()
	c.BaseController.Init(ctx, controllerName, actionName, app)
}

// PrefixIcon 管理界面一级侧栏图标（https://www.layui.com/doc/element/icon.html）
func (c *ManageBaseController) PrefixIcon() string {
	return ""
}

// AdminIcon 管理界面二级侧栏图标（https://www.layui.com/doc/element/icon.html）
func (c *ManageBaseController) AdminIcon() string {
	return ""
}

// List 管理后台列表模板渲染
func (c *ManageBaseController) makeListPK(items *[]DisplayItem) {
	for _, item := range *items {
		if item.PK == "true" {
			c.Data["pkField"] = item.Field
		}
	}
	if c.Data["pkField"] == "" {
		c.Data["pkField"] = (*items)[0].Field
	}
}

// Struct2Map 数据结构体转 map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Invalid {
		return nil
	}

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		switch reflect.TypeOf(field.Interface()).Kind() {
		case reflect.Struct:
			if t.Field(i).Type.Name() == "Time" {
				data[t.Field(i).Name] = v.Field(i).Interface()
			} else {
				data[t.Field(i).Name] = Struct2Map(v.Field(i))
			}
		case reflect.Slice:
			s := reflect.ValueOf(v.Field(i).Interface())
			values := []map[string]interface{}{}
			for _i := 0; _i < s.Len(); _i++ {
				values = append(values, Struct2Map(s.Index(_i).Interface()))
			}
			data[t.Field(i).Name] = values
		default:
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return data
}

// DeepCopy 深拷贝结构体
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
