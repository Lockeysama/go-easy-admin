package geacontrollers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

// 消息码
const (
	MSG_OK  = 0
	MSG_ERR = -1
)

// AjaxMsg ajax返回
func (c *GEAdminBaseController) AjaxMsg(msg interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["message"] = msg
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

// AjaxData ajax返回
func (c *GEAdminBaseController) AjaxData(data interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["data"] = data
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

// AjaxList ajax返回 列表
func (c *GEAdminBaseController) AjaxList(msg interface{}, msgNo int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgNo
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

// AjaxAdd 新增数据
func (c *GEAdminBaseController) AjaxAdd() {
	items := c.DisplayItems(c.Model)
	params := c.parser(items)
	if params == nil {
		c.AjaxMsg("失败", MSG_ERR)
		return
	}
	c.makeListPK(items)

	r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
	var (
		b   []uint8
		err error
	)
	if b, err = json.Marshal(params); err != nil {
		c.AjaxMsg(err.Error(), MSG_ERR)
		return
	}
	if err = json.Unmarshal(b, r); err != nil {
		c.AjaxMsg(err.Error(), MSG_ERR)
		return
	}
	if _, err = c.GEADataBaseInsert(r); err != nil {
		c.AjaxMsg(err.Error(), MSG_ERR)
		return
	}
	for _, item := range *items {
		switch item.DBType {
		case "M2M":
			c.GEADataM2MUpdate(r, item.Field, params[item.Field].([]interface{}), "ADD")
		}
	}

	c.AjaxMsg("成功", MSG_OK)
}

// AjaxUpdate 修改数据
func (c *GEAdminBaseController) AjaxUpdate() {
	items := c.DisplayItems(c.Model)
	params := c.parser(items)
	if params == nil {
		c.AjaxMsg("失败", MSG_ERR)
		return
	}
	c.makeListPK(items)

	m2m := map[string][]interface{}{}
	for _, item := range *items {
		switch item.DBType {
		case DisplayType.M2M:
			if values, ok := params[item.Field]; ok {
				m2m[item.Field] = values.([]interface{})
				delete(params, item.Field)
			}
		}
	}

	if pk, ok := params[c.GetData()["pkField"].(string)]; ok && pk != nil {
		r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
		if _, err := c.GEADataBaseUpdate(
			r.(geamodels.Model),
			map[string]interface{}{c.GetData()["pkField"].(string): pk},
			params,
		); err != nil {
			c.AjaxMsg(err.Error(), MSG_ERR)
			return
		}
		if len(m2m) > 0 {
			row := c.GEADataBaseQueryRow(
				r,
				map[string]interface{}{c.GetData()["pkField"].(string): pk},
				false,
			)
			for name, values := range m2m {
				c.GEADataM2MUpdate(row.(geamodels.Model), name, values, "UPDATE")
			}
		}
	}

	c.AjaxMsg("成功", MSG_OK)
}

// AjaxSave 修改数据
func (c *GEAdminBaseController) parser(displayItems *[]DisplayItem) map[string]interface{} {
	params := map[string]interface{}{}
	for _, item := range *displayItems {
		switch item.DBType {
		case DisplayType.Char, DisplayType.File:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if value[0] == "<nil>" {
					params[item.Field] = nil
				} else {
					params[item.Field] = value[0]
				}
			}
		case DisplayType.Number:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if i, err := strconv.Atoi(value[0]); err == nil {
					params[item.Field] = i
				} else {
					if (value[0] == "<nil>" || value[0] == "") && item.PK != "true" {
						params[item.Field] = nil
					}
				}
			}
		case DisplayType.ForeignKey, DisplayType.O2O:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if i, err := strconv.Atoi(value[0]); err == nil {
					v := reflect.New(reflect.TypeOf(item.Model).Elem()).Elem()
					v.FieldByName(item.Index).SetInt(int64(i))
					params[item.Field] = v.Interface()
				} else {
					if (value[0] == "<nil>" || value[0] == "") && item.PK != "true" {
						params[item.Field] = nil
					}
				}
			}
		case DisplayType.M2M:
			if value, ok := c.RequestForm()[item.Field+"[]"]; ok && len(value) > 0 {
				m2m := []interface{}{}
				m2mDisplayItems := c.DisplayItems(item.Model)
			m2mDisplayItemsLoop:
				for _, m2mItem := range *m2mDisplayItems {
					switch m2mItem.Field {
					case item.Index:
						switch m2mItem.DBType {
						case DisplayType.Number:
							for _, v := range value {
								if i, err := strconv.Atoi(v); err == nil {
									v := reflect.New(reflect.TypeOf(item.Model).Elem()).Elem()
									v.FieldByName(item.Index).SetInt(int64(i))
									m2m = append(m2m, v.Interface())
								}
							}
						default:
							for _, _value := range value {
								v := reflect.New(reflect.TypeOf(item.Model).Elem()).Elem()
								v.FieldByName(item.Index).Set(reflect.ValueOf(_value))
								m2m = append(m2m, v.Interface())
							}
						}
						break m2mDisplayItemsLoop
					}
				}
				params[item.Field] = m2m
			}
		case DisplayType.Bool:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				params[item.Field] = value[0] == "on"
			}
		case DisplayType.Datetime:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if item.DataType == DisplayType.Time {
					params[item.Field], _ = time.Parse("2006-01-02 15:04:05", value[0])
				} else {
					timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", value[0], time.Local)
					params[item.Field] = timestamp.Unix()
				}
			}
		case DisplayType.Time:
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if item.DataType == DisplayType.Time {
					params[item.Field], _ = time.Parse("2006-01-02 15:04:05", value[0])
				} else {
					if ts, err := strconv.Atoi(value[0]); err != nil {
						c.AjaxMsg(err.Error(), 400)
						return nil
					} else {
						params[item.Field] = int64(ts)
					}
				}
			}
		}
		params[item.Field] = DefaultValueMake(params[item.Field], &item)

	}
	return params
}

// AjaxDel 删除数据
func (c *GEAdminBaseController) AjaxDelete() {
	items := c.DisplayItems(c.Model)
	c.makeListPK(items)
	field := c.GetData()["pkField"].(string)
	value := c.RequestQuery(field)
	params := map[string]interface{}{field: value}
	for key := range params {
		switch params[key].(type) {
		case []string:
		case string:
			if _, err := c.GEADataBaseDelete(
				c.Model, map[string]interface{}{key: params[key]},
			); err != nil {
				c.AjaxMsg(err.Error(), MSG_ERR)
				return
			}
		}
	}

	c.AjaxMsg("成功", MSG_OK)
}
