package basecontrollers

import (
	"reflect"
)

// Edit 管理后台编辑模板渲染
func (c *ManageBaseController) Edit() {
	c.Data["pageTitle"] = "编辑"

	gp := c.DisplayItems(c.Model)
	c.makeListPK(gp)

	field := c.Data["pkField"].(string)
	value := c.Ctx.Input.Query(field)

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
