package geacontrollers

import (
	"fmt"
	"strconv"
	"time"
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

// AjaxSave 修改数据
func (c *GEAdminBaseController) AjaxSave() {
	params := map[string]interface{}{}
	displayItems := c.DisplayItems(c.Model)
	c.makeListPK(displayItems)
	for _, item := range *displayItems {
		switch item.DBType {
		case "Char", "File":
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if value[0] == "<nil>" {
					params[item.Field] = nil
				} else {
					params[item.Field] = value[0]
				}
			}
		case "Number", "ForeignKey", "O2O":
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if i, err := strconv.Atoi(value[0]); err == nil {
					params[item.Field] = i
				} else {
					if (value[0] == "<nil>" || value[0] == "") && item.PK != "true" {
						params[item.Field] = nil
					}
				}
			}
		case "M2M":
			if value, ok := c.RequestForm()[item.Field+"[]"]; ok && len(value) > 0 {
				m2m := []interface{}{}
				m2mDisplayItems := c.DisplayItems(item.Model)
			m2mDisplayItemsLoop:
				for _, m2mItem := range *m2mDisplayItems {
					switch m2mItem.Field {
					case item.Index:
						switch m2mItem.DBType {
						case "Number":
							for _, v := range value {
								if i, err := strconv.Atoi(v); err == nil {
									m2m = append(m2m, i)
								}
							}
						default:
							for _, v := range value {
								m2m = append(m2m, v)
							}
						}
						break m2mDisplayItemsLoop
					}
				}
				params[item.Field] = m2m
			}
		case "Bool":
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				params[item.Field] = value[0] == "on"
			}
		case "Datetime":
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if item.DataType == "Time" {
					params[item.Field], _ = time.Parse("2006-01-02 15:04:05", value[0])
				} else {
					timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", value[0], time.Local)
					params[item.Field] = timestamp.Unix()
				}
			}
		case "Time":
			if value, ok := c.RequestForm()[item.Field]; ok && len(value) > 0 {
				if item.DataType == "Time" {
					params[item.Field], _ = time.Parse("2006-01-02 15:04:05", value[0])
				} else {
					if ts, err := strconv.Atoi(value[0]); err != nil {
						c.AjaxMsg(err.Error(), 400)
						return
					} else {
						params[item.Field] = int64(ts)
					}
				}
			}
		}

	}
	items := c.DisplayItems(c.Model)
	c.makeListPK(items)
	if pk, ok := params[c.GetData()["pkField"].(string)]; ok && pk != nil {
		err := c.updateData(pk, displayItems, &params)
		if err != nil {
			c.AjaxMsg(err.Error(), MSG_ERR)
			return
		}
	} else {
		err := c.addData(displayItems, &params)
		if err != nil {
			c.AjaxMsg(err.Error(), MSG_ERR)
			return
		}
	}
	fmt.Println(params)
	c.AjaxMsg("成功", MSG_OK)
}

// AjaxDel 删除数据
func (c *GEAdminBaseController) AjaxDel() {
	items := c.DisplayItems(c.Model)
	c.makeListPK(items)
	field := c.GetData()["pkField"].(string)
	value := c.RequestQuery(field)
	params := map[string]interface{}{field: value}
	for key := range params {
		switch params[key].(type) {
		case []string:
		case string:
			if _, err := c.GEADataBaseDelete(c.Model, map[string]interface{}{key: params[key]}); err != nil {
				c.AjaxMsg(err.Error(), MSG_ERR)
				return
			}
		}
	}

	c.AjaxMsg("成功", MSG_OK)
}

// RequestError API 请求错误
func (c *GEAdminBaseController) RequestError(code int, msg ...string) {
	errMsg := ""
	for _, m := range msg {
		errMsg += (m + ". ")
	}
	if errMsg == "" {
		errMsg = "请求错误"
	}
	c.CustomAbort(code, errMsg)
}
