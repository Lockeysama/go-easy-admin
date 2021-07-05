package basecontrollers

import (
	"reflect"
)

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
