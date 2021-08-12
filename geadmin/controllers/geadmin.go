package geacontrollers

import (
	"bytes"
	"strconv"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"

	"encoding/gob"
	"encoding/json"
	"reflect"

	"github.com/beego/beego/v2/client/orm"
)

// ManageBaseController 控制器管理基类
type ManageBaseController struct {
	BaseController
	Model     geamodels.Model
	PageTitle string
}

// // Init 初始化
// func (c *ManageBaseController) Init(ctx interface{}, controllerName, actionName string, app interface{}) {
// 	c.Instance = app.(ControllerRolePolicy)
// 	c.Model = c.Instance.DBModel()
// 	c.BaseController.Init(ctx, controllerName, actionName, app)
// }

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

// Add 管理后台新增模板渲染
func (c *ManageBaseController) Add() {
	c.Data["pageTitle"] = "新增"
	items := c.DisplayItems(c.Model)
	c.makeListPK(items)
	linkItemsMap := map[string][]map[string]interface{}{}
	for _, item := range *items {
		if item.DBType == "M2M" || item.DBType == "O2O" || item.DBType == "ForeignKey" {
			linkItems := c.QueryList(item.Model, 0, 0, nil, nil, false)
			linkValue := reflect.ValueOf(linkItems).Elem()
			linksMap := []map[string]interface{}{}
			for i := 0; i < linkValue.Len(); i++ {
				linkMap := map[string]interface{}{}
				sm := Struct2Map(linkValue.Index(i).Elem().Interface())
				linkMap[item.ShowField] = sm[item.ShowField]
				linkMap[item.Index] = sm[item.Index]
				linksMap = append(linksMap, linkMap)
			}

			linkItemsMap[item.Field] = linksMap
		}
	}
	c.Data["linkItems"] = linkItemsMap
	c.Data["display"] = items
	c.Layout = "public/layout.html"
	c.TplName = "public/add.html"
}

// List 管理后台列表模板渲染
func (c *ManageBaseController) List() {
	c.Data["pageTitle"] = "列表"
	items := c.DisplayItems(c.Model)
	c.Data["display"] = items
	c.makeListPK(items)
	c.Layout = "public/layout.html"
	c.TplName = "public/list.html"
}

// Table 获取管理后台列表数据
func (c *ManageBaseController) Table() {
	// 列表
	page, err := strconv.Atoi(c.Ctx().InputParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Ctx().InputParam("limit"))
	if err != nil {
		limit = 30
	}

	var lists interface{}

	orderByStr := c.Ctx().InputQuery("order_by")
	listOrderBy := make(map[string]string)
	if orderByStr != "" {
		if err := json.Unmarshal([]byte(orderByStr), &listOrderBy); err != nil {
			c.AjaxMsg("order exception", MSG_ERR)
			return
		}
	}

	query := c.Ctx().InputQuery("query")
	listFilter := make(map[string]interface{})

	if query != "" {
		listFilter = c.ListFilter()
		lists = c.QueryList(c.Model, page, limit, listFilter, listOrderBy, true)
	} else {
		lists = c.QueryList(c.Model, page, limit, nil, listOrderBy, true)
	}

	b, _ := json.Marshal(lists)
	var resultMap = new([]map[string]interface{})
	json.Unmarshal(b, resultMap)

	c.PageSize = limit
	listsValue := reflect.ValueOf(lists).Elem()
	list := make([]map[string]interface{}, listsValue.Len())
	for i := 0; i < listsValue.Len(); i++ {
		x := Struct2MapWithHTML(
			&(*resultMap)[i], c.DisplayItems(c.Model),
		)
		list[i] = x
	}
	if query != "" {
		if count, err := c.QueryCount(c.Model, listFilter); err != nil {
			c.AjaxMsg("查询失败", MSG_ERR)
		} else {
			c.AjaxList("成功", MSG_OK, count, list)
		}
	} else {
		if count, err := c.QueryCount(c.Model, nil); err != nil {
			c.AjaxMsg("查询失败", MSG_ERR)
		} else {
			c.AjaxList("成功", MSG_OK, count, list)
		}
	}
}

func (c *ManageBaseController) ListFilter() map[string]interface{} {
	queryStr := c.Ctx().InputQuery("query")
	if queryStr == "" {
		c.AjaxMsg("请输入查询条件", MSG_ERR)
		return nil
	}

	query := new([]map[string]interface{})
	if err := json.Unmarshal([]byte(queryStr), query); err != nil {
		c.AjaxMsg("查询条件异常", MSG_ERR)
		return nil
	}

	if len(*query) == 0 {
		c.AjaxMsg("查询条件异常", MSG_ERR)
		return nil
	}

	filters := make(map[string]interface{})

	for _, q := range *query {
		ok := false
		field := ""
		if field, ok = q["field"].(string); !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}
		expression := ""
		if expression, ok = q["exp"].(string); !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}
		var value interface{}
		if value, ok = q["value"]; !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}

		switch expression {
		case "eq":
			filters[field] = value
		case "ne":
			filters[field+"__iexact"] = value
		case "gt":
			filters[field+"__gt"] = value
		case "lt":
			filters[field+"__lt"] = value
		case "gte":
			filters[field+"__gte"] = value
		case "lte":
			filters[field+"__lte"] = value
		case "is_contains":
			filters[field+"__contains"] = value
		case "not_contains":
			filters[field+"__icontains"] = value
		}
	}
	return filters
}

// Edit 管理后台编辑模板渲染
func (c *ManageBaseController) Edit() {
	c.Data["pageTitle"] = "编辑"

	gp := c.DisplayItems(c.Model)
	c.makeListPK(gp)

	field := c.Data["pkField"].(string)
	value := c.Ctx().InputQuery(field)

	filters := map[string]interface{}{field: value}
	r := c.QueryRow(c.Model, filters, true)

	if r == nil {
		c.AjaxMsg("data exception", MSG_ERR)
		return
	} else {
		v := reflect.ValueOf(r).Elem()
		linkItemsMap := map[string][]map[string]interface{}{}
		for i, item := range *gp {
			if item.DBType == "Datetime" {
				if v.FieldByName(item.Field).Type().Name() == "Time" {
					(*gp)[i].Value = v.FieldByName(item.Field).Interface()
				} else {
					(*gp)[i].Value = v.FieldByName(item.Field).Int() * 1000
				}
			} else if item.DBType == "M2M" || item.DBType == "O2O" || item.DBType == "ForeignKey" {
				itemValuesMap := []map[string]interface{}{}
				switch v.FieldByName(item.Field).Type().Kind() {
				case reflect.Slice:
					itemValues := v.FieldByName(item.Field)
					for _i := 0; _i < itemValues.Len(); _i++ {
						itemValuesMap = append(itemValuesMap, Struct2Map(itemValues.Index(_i).Elem().Interface()))
					}
				default:
					itemValues := v.FieldByName(item.Field)
					itemValuesMap = append(itemValuesMap, Struct2Map(itemValues.Interface()))
				}
				linkItems := c.QueryList(item.Model, 0, 0, nil, nil, false)
				linkValue := reflect.ValueOf(linkItems).Elem()
				linksMap := []map[string]interface{}{}
				for i := 0; i < linkValue.Len(); i++ {
					linkMap := map[string]interface{}{}
					sm := Struct2Map(linkValue.Index(i).Elem().Interface())
					linkMap[item.ShowField] = sm[item.ShowField]
					linkMap[item.Index] = sm[item.Index]
					for _, ivm := range itemValuesMap {
						if linkMap[item.Index] == ivm[item.Index] {
							linkMap["checked"] = true
							break
						}
						linkMap["checked"] = false
					}
					linksMap = append(linksMap, linkMap)
				}

				linkItemsMap[item.Field] = linksMap
			} else {
				(*gp)[i].Value = v.FieldByName(item.Field).Interface()
			}
		}
		c.Data["linkItems"] = linkItemsMap
	}

	c.Data["display"] = gp
	c.Layout = "public/layout.html"
	c.TplName = "public/edit.html"
}

// Detail 管理后台编辑模板渲染
func (c *ManageBaseController) Detail() {
	c.Data["pageTitle"] = " 详情"

	gp := c.DisplayItems(c.Model)
	c.makeListPK(gp)

	var detailDisplayItems *[]DisplayItem
	for _, item := range *gp {
		field := c.Ctx().InputQuery("field")
		if item.Field == field {
			index, _ := strconv.Atoi(c.Ctx().InputParam(item.Index))
			filters := map[string]interface{}{item.Index: index}
			r := c.QueryRow(item.Model, filters, false)
			if detailDisplayItems = c.DisplayItems(r.(geamodels.Model)); len(*detailDisplayItems) < 1 {
				panic("DisplayItems Exception")
			}

			if r == nil {
				c.AjaxMsg("data exception", MSG_ERR)
				return
			} else {
				v := reflect.ValueOf(r).Elem()
				linkItemsMap := map[string][]map[string]interface{}{}
				for i, detailItem := range *detailDisplayItems {
					if detailItem.DBType == "Datetime" {
						if detailItem.DataType == "Time" {
							(*detailDisplayItems)[i].Value = v.FieldByName(detailItem.Field).Interface()
						} else {
							(*detailDisplayItems)[i].Value = v.FieldByName(detailItem.Field).Int() * 1000
						}
					} else if detailItem.DBType == "M2M" || detailItem.DBType == "O2O" || detailItem.DBType == "ForeignKey" {
						itemValuesMap := []map[string]interface{}{}
						switch v.FieldByName(detailItem.Field).Type().Kind() {
						case reflect.Slice:
							orm.NewOrm().LoadRelated(r, detailItem.Field)
							itemValues := v.FieldByName(detailItem.Field)
							for _i := 0; _i < itemValues.Len(); _i++ {
								itemValuesMap = append(itemValuesMap, Struct2Map(itemValues.Index(_i).Elem().Interface()))
							}
						default:
							itemValues := v.FieldByName(detailItem.Field)
							itemValuesMap = append(itemValuesMap, Struct2Map(itemValues.Interface()))
						}

						linkItems := c.QueryList(detailItem.Model, 0, 0, nil, nil, false)
						linkValue := reflect.ValueOf(linkItems).Elem()
						linksMap := []map[string]interface{}{}
						for i := 0; i < linkValue.Len(); i++ {
							linkMap := map[string]interface{}{}
							sm := Struct2Map(linkValue.Index(i).Elem().Interface())
							linkMap[detailItem.ShowField] = sm[detailItem.ShowField]
							linkMap[detailItem.Index] = sm[detailItem.Index]
							for _, ivm := range itemValuesMap {
								if linkMap[detailItem.Index] == ivm[detailItem.Index] {
									linkMap["checked"] = true
									break
								}
								linkMap["checked"] = false
							}
							linksMap = append(linksMap, linkMap)
						}

						linkItemsMap[detailItem.Field] = linksMap
					} else {
						(*detailDisplayItems)[i].Value = v.FieldByName(detailItem.Field).Interface()
					}
				}
				c.Data["linkItems"] = linkItemsMap
			}
			// break
		}
	}

	c.Data["display"] = detailDisplayItems
	c.Layout = "public/layout.html"
	c.TplName = "public/detail.html"
}
