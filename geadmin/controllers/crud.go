package geacontrollers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"

	"github.com/beego/beego/v2/client/orm"
)

// QueryCount 修改数据
func (c *GEAManageBaseController) QueryCount(model geamodels.Model, filters map[string]interface{}) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(model)
	for field := range filters {
		qs = qs.Filter(field, filters[field])
	}

	return qs.Count()
}

// QueryList 修改数据
func (c *GEAManageBaseController) QueryListWithRawFilter(
	model geamodels.Model, page int, limit int, filters string, loadRel bool,
) interface{} {
	lists := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()

	o := orm.NewOrm()
	qs := o.QueryTable(model)
	if limit > 0 {
		qs = qs.Limit(limit, (page-1)*limit)
	}
	if filters != "" {
		qs = qs.FilterRaw(filters, "")
	}
	qs.All(lists)

	if loadRel {
		fksList := make(map[interface{}]reflect.Value)
		vLists := reflect.ValueOf(lists).Elem()
		for _, item := range *c.DisplayItems(model) {
			switch item.DBType {
			case "M2M":
				for i := 0; i < vLists.Len(); i++ {
					o.LoadRelated(vLists.Index(i).Interface(), item.Field)
				}
			case "O2O", "ForeignKey":
				if _, ok := fksList[item.Model]; !ok {
					fksList[item.Model] = reflect.ValueOf(c.QueryList(item.Model, 0, 0, nil, nil, false)).Elem()
				}
				for i := 0; i < vLists.Len(); i++ {
					vRow := vLists.Index(i).Elem()
					if vRow.FieldByName(item.Field).IsNil() {
						continue
					}
					vRowIndex := vRow.FieldByName(item.Field).Elem().FieldByName(item.Index).Interface()
					vLists.Index(i).Elem().FieldByName(item.Field).
						Set(reflect.Zero(vLists.Index(i).Elem().FieldByName(item.Field).Type()))
					for j := 0; j < fksList[item.Model].Len(); j++ {
						if vRowIndex == fksList[item.Model].Index(j).Elem().FieldByName(item.Index).Interface() {
							vLists.Index(i).Elem().FieldByName(item.Field).Set(fksList[item.Model].Index(j))
							break
						}
					}
				}
			}
		}
	}
	return lists
}

// QueryList 修改数据
func (c *GEAManageBaseController) QueryList(
	model geamodels.Model,
	page int, limit int,
	filters map[string]interface{},
	order map[string]string,
	loadRel bool,
) interface{} {
	lists := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()

	o := orm.NewOrm()
	qs := o.QueryTable(model)

	for field := range filters {
		qs = qs.Filter(field, filters[field])
	}

	for field := range order {
		orderType := order[field]
		if orderType == "desc" {
			qs = qs.OrderBy("-" + field)
		} else {
			qs = qs.OrderBy(field)
		}
	}

	if limit > 0 {
		qs = qs.Limit(limit, (page-1)*limit)
	}

	count, _ := qs.All(lists)
	fmt.Println(count)

	if loadRel {
		fksList := make(map[interface{}]reflect.Value)
		vLists := reflect.ValueOf(lists).Elem()
		for _, item := range *c.DisplayItems(model) {
			switch item.DBType {
			case "M2M":
				for i := 0; i < vLists.Len(); i++ {
					vLists.Index(i).Interface().(geamodels.Model).LoadM2M()
				}
			case "O2O", "ForeignKey":
				if _, ok := fksList[item.Model]; !ok {
					fksList[item.Model] = reflect.ValueOf(
						c.QueryList(item.Model, 0, 0, nil, nil, false),
					).Elem()
				}
				for i := 0; i < vLists.Len(); i++ {
					vRow := vLists.Index(i).Elem()
					if vRow.FieldByName(item.Field).IsNil() {
						continue
					}
					vRowIndex := vRow.FieldByName(item.Field).Elem().FieldByName(item.Index).Interface()
					vLists.Index(i).Elem().FieldByName(item.Field).
						Set(reflect.Zero(vLists.Index(i).Elem().FieldByName(item.Field).Type()))
					for j := 0; j < fksList[item.Model].Len(); j++ {
						if vRowIndex == fksList[item.Model].Index(j).Elem().FieldByName(item.Index).Interface() {
							vLists.Index(i).Elem().FieldByName(item.Field).Set(fksList[item.Model].Index(j))
							break
						}
					}
				}
			}
		}
	}
	return lists
}

// QueryRow 修改数据
func (c *GEAManageBaseController) QueryRow(
	model geamodels.Model, filters map[string]interface{}, loadRel bool,
) interface{} {
	row := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	o := orm.NewOrm()
	qs := o.QueryTable(model)
	for field := range filters {
		qs = qs.Filter(field, filters[field])
	}

	qs.One(row)

	if loadRel {
		fksList := make(map[interface{}]reflect.Value)
		for _, item := range *c.DisplayItems(model) {
			switch item.DBType {
			case "M2M":
				row.(geamodels.Model).LoadM2M()
			case "O2O", "ForeignKey":
				if _, ok := fksList[item.Model]; !ok {
					fksList[item.Model] = reflect.ValueOf(
						c.QueryList(item.Model, 0, 0, nil, nil, false),
					).Elem()
				}
				vRow := reflect.ValueOf(row).Elem()
				if vRow.FieldByName(item.Field).IsNil() {
					continue
				}
				vRowIndex := vRow.FieldByName(item.Field).Elem().FieldByName(item.Index).Interface()
				for j := 0; j < fksList[item.Model].Len(); j++ {
					if vRowIndex == fksList[item.Model].Index(j).Elem().FieldByName(item.Index).Interface() {
						vRow.FieldByName(item.Field).Set(fksList[item.Model].Index(j))
						break
					}
				}
			}
		}
	}
	return row
}

// AjaxSave 修改数据
func (c *GEAManageBaseController) AjaxSave() {
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

func (c *GEAManageBaseController) addData(displayItems *[]DisplayItem, params *map[string]interface{}) error {
	o := orm.NewOrm()
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
		var (
			b   []uint8
			err error
		)
		if b, err = json.Marshal(*params); err != nil {
			fmt.Println(err.Error())
		}
		json.Unmarshal(b, r)

		if _, err := txOrm.Insert(r); err != nil {
			return err
		}
		for _, item := range *displayItems {
			switch item.DBType {
			case "M2M":
				m2m := txOrm.QueryM2M(r, item.Field)
				if indexes, ok := (*params)[item.Field]; ok {
					if len(indexes.([]interface{})) > 0 {
						m2m.Add(indexes.([]interface{}))
					}
				}
			case "ForeignKey", "O2O":
				if index, ok := (*params)[item.Field]; ok {
					rIndex := reflect.ValueOf(r).Elem().FieldByName(c.GetData()["pkField"].(string)).Interface()
					qs := txOrm.QueryTable(c.Model).Filter(c.GetData()["pkField"].(string), rIndex)
					// target := reflect.New(reflect.TypeOf(item.Model).Elem()).Interface()
					// txOrm.QueryTable(item.Model).Filter(item.Index, index).One(target)
					if _, err := qs.Update(map[string]interface{}{item.Field: index}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	return err
}

func (c *GEAManageBaseController) updateData(pk interface{}, displayItems *[]DisplayItem, params *map[string]interface{}) error {
	o := orm.NewOrm()
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		qs := txOrm.QueryTable(c.Model).Filter(c.GetData()["pkField"].(string), pk)
		r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
		qs.One(r)
		for _, item := range *displayItems {
			switch item.DBType {
			case "M2M":
				txOrm.LoadRelated(r, item.Field)
				old := []interface{}{}
				m2mRowsValue := reflect.ValueOf(r).Elem().FieldByName(item.Field)
				for i := 0; i < m2mRowsValue.Len(); i++ {
					old = append(old, m2mRowsValue.Index(i).Elem().FieldByName(item.Index).Interface())
				}
				m2m := txOrm.QueryM2M(r, item.Field)
				if len(old) > 0 {
					m2m.Remove(old)
				}
				if indexes, ok := (*params)[item.Field]; ok {
					if len(indexes.([]interface{})) > 0 {
						m2m.Add(indexes.([]interface{}))
					}
				}
				delete(*params, item.Field)
			}
		}
		if _, err := qs.Update(*params); err != nil {
			return err
		}
		return nil
	})
	return err
}

// AjaxDel 删除数据
func (c *GEAManageBaseController) AjaxDel() {
	items := c.DisplayItems(c.Model)
	c.makeListPK(items)
	field := c.GetData()["pkField"].(string)
	value := c.RequestQuery(field)
	params := map[string]interface{}{field: value}
	for key := range params {
		switch params[key].(type) {
		case []string:
		case string:
			if _, err := orm.NewOrm().
				QueryTable(c.Model).
				Filter(key, params[key]).
				Delete(); err != nil {
				c.AjaxMsg(err.Error(), MSG_ERR)
				return
			}
		}
	}

	c.AjaxMsg("成功", MSG_OK)
}
