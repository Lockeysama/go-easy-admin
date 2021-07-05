package basecontrollers

import (
	basemodels "TDCS/models/base"
	"reflect"

	"github.com/beego/beego/v2/client/orm"
)

// Detail 管理后台编辑模板渲染
func (c *ManageBaseController) Detail() {
	c.Data["pageTitle"] = " 详情"

	gp := c.DisplayItems(c.Model)
	c.makeListPK(gp)

	var detailDisplayItems *[]DisplayItem
	for _, item := range *gp {
		field := c.GetString("field", "")
		if item.Field == field {
			index, _ := c.GetInt(item.Index, 0)
			filters := map[string]interface{}{item.Index: index}
			r := c.QueryRow(item.Model, filters, false)
			if detailDisplayItems = c.DisplayItems(r.(basemodels.Model)); len(*detailDisplayItems) < 1 {
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
