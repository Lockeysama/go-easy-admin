package geacontrollers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// ModelDataCleaning 模型数据清洗
func (c *GEABaseController) ModelDataCleaning(data interface{}) interface{} {
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		data = reflect.ValueOf(data).Elem().Interface()
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		result := []interface{}{}
		list := reflect.ValueOf(data)
		for i := 0; i < list.Len(); i++ {
			row := list.Index(i)
			result = append(result, RowDataCleaning(row))
		}
		return result
	case reflect.Struct:
		return RowDataCleaning(reflect.ValueOf(data))
	default:
		return []interface{}{}
	}
}

// RowDataCleaning 行数据清洗
func RowDataCleaning(row reflect.Value) interface{} {
	if row.Type().Kind() == reflect.Ptr {
		row = row.Elem()
	}
	result := make(map[string]interface{})
	for i := 0; i < row.NumField(); i++ {
		v := row.Field(i)
		if row.Type().Field(i).Tag.Get("json") == "-" {
			continue
		}
		switch v.Type().Kind() {
		case reflect.Slice:
			ids := []interface{}{}
			for mi := 0; mi < v.Len(); mi++ {
				mv := v.Index(mi)
				if mv.Type().Kind() == reflect.Ptr {
					ids = append(ids, mv.Elem().FieldByName("ID").Interface()) // TODO ID 换主键
				} else {
					ids = append(ids, mv.FieldByName("ID").Interface())
				}
			}
			if len(ids) > 0 {
				result[row.Type().Field(i).Name] = ids
			}
		case reflect.Ptr:
			if !v.IsNil() {
				if v.Elem().Type().Kind() == reflect.Struct {
					result[row.Type().Field(i).Name] = v.Elem().FieldByName("ID").Interface()
				} else {
					result[row.Type().Field(i).Name] = v.Elem().Interface()
				}
			}
		case reflect.Struct:
			if v.Type().Name() != "Time" {
				result[row.Type().Field(i).Name] = v.FieldByName("ID").Interface()
			} else {
				result[row.Type().Field(i).Name] = v.Interface()
			}
		default:
			result[row.Type().Field(i).Name] = v.Interface()
		}
	}
	return result
}

// RequestID API Get 请求 ID
func (c *GEABaseController) RequestID() int64 {
	id := c.Ctx().InputParam(":id")
	if id, err := strconv.Atoi(id); err != nil {
		c.APIRequestError(400, err.Error())
	} else {
		return int64(id)
	}
	return -1
}

// APIRequestError API 请求错误
func (c *GEABaseController) APIRequestError(code int, msg ...string) {
	errMsg := ""
	for _, m := range msg {
		errMsg += (m + ". ")
	}
	if errMsg == "" {
		errMsg = "请求错误"
	}
	c.CustomAbort(code, errMsg)
}

// Post 创建数据
func (c *GEAManageBaseController) Post() {
	rowMap := make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx().InputRequestBody(), &rowMap); err != nil {
		c.APIRequestError(400, err.Error())
		return
	}
	displayItems := c.DisplayItems(c.Model)
	for _, item := range *displayItems {
		if _, ok := rowMap[item.Field]; !ok {
			continue
		}
		if item.PK == "true" {
			delete(rowMap, item.Field)
		}
		switch item.DBType {
		case "ForeignKey", "O2O":
			rowMap[item.Field] = map[string]interface{}{"ID": rowMap[item.Field]}
		case "M2M":
			m2m := []map[string]interface{}{}
			for _, v := range rowMap[item.Field].([]int64) {
				m2m = append(m2m, map[string]interface{}{"ID": v})
			}
			rowMap[item.Field] = m2m
		case "Datetime":
			if datetime, err := time.Parse(time.RFC3339, rowMap[item.Field].(string)); err != nil {
				c.APIRequestError(400, err.Error())
				return
			} else {
				if !datetime.IsZero() {
					rowMap[item.Field] = datetime
				} else {
					delete(rowMap, item.Field)
				}
			}
		}
	}
	if realBody, err := json.Marshal(rowMap); err != nil {
		c.APIRequestError(400, err.Error())
		return
	} else {
		r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
		if err := json.Unmarshal(realBody, r); err != nil {
			c.APIRequestError(400, err.Error())
			return
		}

		o := orm.NewOrm()
		if id, err := o.Insert(r); err != nil {
			c.APIRequestError(400, err.Error())
			return
		} else {
			qs := o.QueryTable(c.Model).Filter("id", id)
			if err := qs.One(r); err != nil {
				c.APIRequestError(400, err.Error())
				return
			}
		}

		for _, item := range *displayItems {
			if _, ok := rowMap[item.Field]; !ok {
				continue
			}
			switch item.DBType {
			case "M2M":
				o.QueryM2M(r, item.Field).Add(rowMap[item.Field])
				o.LoadRelated(r, item.Field)
			}
		}
		c.SetData("json", c.ModelDataCleaning(r))

		c.ServeJSON()
	}
}

// GetAll 获取所有数据
func (c *GEAManageBaseController) GetAll(filters map[string]interface{}, order string) {
	lists := reflect.New(reflect.SliceOf(reflect.TypeOf(c.Model))).Interface()

	// 列表
	page, err := strconv.Atoi(c.Ctx().InputQuery("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Ctx().InputQuery("limit"))
	if err != nil {
		limit = 0
	}

	o := orm.NewOrm()
	qs := o.QueryTable(c.Model)

	if order != "" {
		qs = qs.OrderBy(order)
	}

	if limit > 0 {
		qs = qs.Limit(limit, (page-1)*limit)
	}

	for k, v := range filters {
		qs = qs.Filter(k, v)
	}

	qs.All(lists)

	displayItems := c.DisplayItems(c.Model)
	for _, item := range *displayItems {
		switch item.DBType {
		case "M2M":
			listsValue := reflect.ValueOf(lists).Elem()
			for i := 0; i < listsValue.Len(); i++ {
				o.LoadRelated(listsValue.Index(i).Interface(), item.Field)
			}
		}
	}

	// c.SetData("json", c.ModelDataCleaning(lists)
	c.SetData("json", c.ModelDataCleaning(lists))
	c.ServeJSON()
}

// Get 获取单条数据
func (c *GEAManageBaseController) Get() {
	row := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()

	o := orm.NewOrm()
	qs := o.QueryTable(c.Model).Filter("id", c.RequestID())
	if err := qs.One(row); err != nil {
		c.APIRequestError(404, err.Error())
		return
	}
	displayItems := c.DisplayItems(c.Model)
	for _, item := range *displayItems {
		switch item.DBType {
		case "M2M":
			o.LoadRelated(row, item.Field)
		}
	}
	c.SetData("json", c.ModelDataCleaning(row))
	c.ServeJSON()
}

// GetWith Filters 获取单条数据
func (c *GEAManageBaseController) GetWith(filters ...map[string]interface{}) {
	row := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()

	o := orm.NewOrm()
	qs := o.QueryTable(c.Model).Filter("id", c.RequestID())
	for _, filter := range filters {
		for k, v := range filter {
			qs = qs.Filter(k, v)
		}
	}
	if err := qs.One(row); err != nil {
		c.APIRequestError(404, err.Error())
		return
	}
	displayItems := c.DisplayItems(c.Model)
	for _, item := range *displayItems {
		switch item.DBType {
		case "M2M":
			o.LoadRelated(row, item.Field)
		}
	}
	c.SetData("json", c.ModelDataCleaning(row))
	c.ServeJSON()
}

// Put 更新数据
func (c *GEAManageBaseController) Put() {
	params := make(map[string]interface{})
	json.Unmarshal(c.Ctx().InputRequestBody(), &params)

	displayItems := c.DisplayItems(c.Model)
	m2m := map[string][]interface{}{}
	for field := range params {
		for _, item := range *displayItems {
			if strings.EqualFold(strings.ToLower(field), strings.ToLower(item.Field)) {
				switch item.DBType {
				case "M2M":
					m2m[item.Field] = params[field].([]interface{})
					delete(params, field)
				case "Datetime":
					if datetime, err := time.Parse(time.RFC3339, params[field].(string)); err != nil {
						c.APIRequestError(400, err.Error())
						return
					} else {
						if !datetime.IsZero() {
							params[field] = datetime
						} else {
							delete(params, field)
						}
					}
				}
			}
		}
	}

	o := orm.NewOrm()
	requestID := c.RequestID()
	qs := o.QueryTable(c.Model).Filter("ID", requestID)

	if len(params) > 0 {
		if _, err := qs.Update(params); err != nil {
			c.APIRequestError(403, err.Error())
			return
		}
	}

	row := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
	if err := qs.One(row); err != nil {
		c.APIRequestError(403, err.Error())
		return
	}

	rowValue := reflect.ValueOf(row)
	for m2mField := range m2m {
		o.LoadRelated(row, m2mField)
		ids := []interface{}{}
		for i := 0; i < rowValue.Elem().FieldByName(m2mField).Len(); i++ {
			id := rowValue.Elem().FieldByName(m2mField).Index(i).Elem().FieldByName("ID").Interface()
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			o.QueryM2M(row, m2mField).Remove(ids)
			o.LoadRelated(row, m2mField)
		}
		if len(m2m[m2mField]) > 0 {
			o.QueryM2M(row, m2mField).Add(m2m[m2mField])
			o.LoadRelated(row, m2mField)
		}
	}

	c.SetData("json", c.ModelDataCleaning(row))
	c.ServeJSON()
}

// PutWith Filters 更新数据
func (c *GEAManageBaseController) PutWith(filters ...map[string]interface{}) {
	params := make(map[string]interface{})
	json.Unmarshal(c.Ctx().InputRequestBody(), &params)

	displayItems := c.DisplayItems(c.Model)
	m2m := map[string][]interface{}{}
	for field := range params {
		for _, item := range *displayItems {
			if strings.EqualFold(strings.ToLower(field), strings.ToLower(item.Field)) {
				switch item.DBType {
				case "M2M":
					m2m[item.Field] = params[field].([]interface{})
					delete(params, field)
				case "Datetime":
					if datetime, err := time.Parse(time.RFC3339, params[field].(string)); err != nil {
						c.APIRequestError(400, err.Error())
						return
					} else {
						if !datetime.IsZero() {
							params[field] = datetime
						} else {
							delete(params, field)
						}
					}
				}
			}
		}
	}

	o := orm.NewOrm()
	qs := o.QueryTable(c.Model).Filter("ID", c.RequestID())
	for _, filter := range filters {
		for k, v := range filter {
			qs = qs.Filter(k, v)
		}
	}

	if len(params) > 0 {
		if _, err := qs.Update(params); err != nil {
			c.APIRequestError(403, err.Error())
			c.ServeJSON()
		}
	}

	row := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
	qs.One(row)

	rowValue := reflect.ValueOf(row)
	for m2mField := range m2m {
		o.LoadRelated(row, m2mField)
		ids := []interface{}{}
		for i := 0; i < rowValue.Elem().FieldByName(m2mField).Len(); i++ {
			id := rowValue.Elem().FieldByName(m2mField).Index(i).Elem().FieldByName("ID").Interface()
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			o.QueryM2M(row, m2mField).Remove(ids)
			o.LoadRelated(row, m2mField)
		}
		if len(m2m[m2mField]) > 0 {
			o.QueryM2M(row, m2mField).Add(m2m[m2mField])
			o.LoadRelated(row, m2mField)
		}
	}

	c.SetData("json", c.ModelDataCleaning(row))
	c.ServeJSON()
}

// Delete 删除数据
func (c *GEAManageBaseController) Delete() {
	qs := orm.NewOrm().QueryTable(c.Model).Filter("id", c.RequestID())
	if _, err := qs.Delete(); err != nil {
		c.AjaxMsg(err.Error(), MSG_ERR)
		return
	}
	c.SetData("json", map[string]int{"code": 0})
	c.ServeJSON()
}

// DeleteWith 删除数据
func (c *GEAManageBaseController) DeleteWith(filters ...map[string]interface{}) {
	qs := orm.NewOrm().QueryTable(c.Model).Filter("id", c.RequestID())
	for _, filter := range filters {
		for k, v := range filter {
			qs = qs.Filter(k, v)
		}
	}
	if _, err := qs.Delete(); err != nil {
		c.AjaxMsg(err.Error(), MSG_ERR)
		return
	}
	c.SetData("json", map[string]int{"code": 0})
	c.ServeJSON()
}
