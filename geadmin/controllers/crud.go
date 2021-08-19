package geacontrollers

import (
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

type GEADataBase interface {
	GEADataBaseCount(model geamodels.Model, filters map[string]interface{}) (int64, error)
	GEADataBaseQueryList(
		model geamodels.Model,
		page int, limit int,
		filters map[string]interface{},
		order map[string]string,
		loadRel bool,
	) interface{}
	GEADataBaseQueryRow(model geamodels.Model, filters map[string]interface{}, loadRel bool) interface{}

	GEADataBaseInsert(model geamodels.Model) (int64, error)
	GEADataBaseUpdate(model geamodels.Model, filters map[string]interface{}, params map[string]interface{}) (int64, error)
	GEADataBaseDelete(model geamodels.Model, filters map[string]interface{}) (int64, error)
}

// QueryCount 修改数据
// func (c *GEAdminBaseController) GEADataBaseCount(
// 	model geamodels.Model, filters map[string]interface{},
// ) (int64, error) {
// 	panic("QueryCount dose not adapt")
// }

// QueryList 修改数据
// func (c *GEAdminBaseController) GEADataBaseQueryList(
// 	model geamodels.Model,
// 	page int, limit int,
// 	filters map[string]interface{},
// 	order map[string]string,
// 	loadRel bool,
// ) interface{} {
// 	panic("QueryList dose not adapt")
// }

// QueryRow 修改数据
// func (c *GEAdminBaseController) GEADataBaseQueryRow(
// 	model geamodels.Model, filters map[string]interface{}, loadRel bool,
// ) interface{} {
// 	panic("QueryRow dose not adapt")
// }

// func (c *GEAdminBaseController) GEADataBaseInsert(model geamodels.Model) (int64, error) {
// 	panic("Insert dose not adapt")
// }

// func (c *GEAdminBaseController) GEADataBaseUpdate(
// 	model geamodels.Model, filters map[string]interface{}, params map[string]interface{},
// ) (int64, error) {
// 	panic("Update dose not adapt")
// }

// func (c *GEAdminBaseController) GEADataBaseDelete(
// 	model geamodels.Model, filters map[string]interface{},
// ) (int64, error) {
// 	panic("Delete dose not adapt")
// }

func (c *GEAdminBaseController) addData(
	displayItems *[]DisplayItem, params *map[string]interface{},
) error {
	panic("Insert dose not adapt")
	// o := orm.NewOrm()
	// err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
	// 	r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
	// 	var (
	// 		b   []uint8
	// 		err error
	// 	)
	// 	if b, err = json.Marshal(*params); err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	json.Unmarshal(b, r)

	// 	if _, err := txOrm.Insert(r); err != nil {
	// 		return err
	// 	}
	// 	for _, item := range *displayItems {
	// 		switch item.DBType {
	// 		case "M2M":
	// 			m2m := txOrm.QueryM2M(r, item.Field)
	// 			if indexes, ok := (*params)[item.Field]; ok {
	// 				if len(indexes.([]interface{})) > 0 {
	// 					m2m.Add(indexes.([]interface{}))
	// 				}
	// 			}
	// 		case "ForeignKey", "O2O":
	// 			if index, ok := (*params)[item.Field]; ok {
	// 				rIndex := reflect.ValueOf(r).Elem().FieldByName(c.GetData()["pkField"].(string)).Interface()
	// 				qs := txOrm.QueryTable(c.Model).Filter(c.GetData()["pkField"].(string), rIndex)
	// 				// target := reflect.New(reflect.TypeOf(item.Model).Elem()).Interface()
	// 				// txOrm.QueryTable(item.Model).Filter(item.Index, index).One(target)
	// 				if _, err := qs.Update(map[string]interface{}{item.Field: index}); err != nil {
	// 					return err
	// 				}
	// 			}
	// 		}
	// 	}
	// 	return nil
	// })
	// return err
}

func (c *GEAdminBaseController) updateData(pk interface{}, displayItems *[]DisplayItem, params *map[string]interface{}) error {
	panic("Update dose not adapt")
	// o := orm.NewOrm()
	// err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
	// 	qs := txOrm.QueryTable(c.Model).Filter(c.GetData()["pkField"].(string), pk)
	// 	r := reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
	// 	qs.One(r)
	// 	for _, item := range *displayItems {
	// 		switch item.DBType {
	// 		case "M2M":
	// 			txOrm.LoadRelated(r, item.Field)
	// 			old := []interface{}{}
	// 			m2mRowsValue := reflect.ValueOf(r).Elem().FieldByName(item.Field)
	// 			for i := 0; i < m2mRowsValue.Len(); i++ {
	// 				old = append(old, m2mRowsValue.Index(i).Elem().FieldByName(item.Index).Interface())
	// 			}
	// 			m2m := txOrm.QueryM2M(r, item.Field)
	// 			if len(old) > 0 {
	// 				m2m.Remove(old)
	// 			}
	// 			if indexes, ok := (*params)[item.Field]; ok {
	// 				if len(indexes.([]interface{})) > 0 {
	// 					m2m.Add(indexes.([]interface{}))
	// 				}
	// 			}
	// 			delete(*params, item.Field)

	// 		}
	// 	}
	// 	if _, err := qs.Update(*params); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	// return err
}
