package basecontrollers

import (
	"fmt"
	"reflect"

	"github.com/beego/beego/v2/client/orm"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func (c *AdaptController) GEADataBaseCount(
	model geamodels.Model, filters map[string]interface{},
) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(model)
	for field := range filters {
		qs = qs.Filter(field, filters[field])
	}

	return qs.Count()
}

func (c *AdaptController) GEADataBaseQueryList(
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
			case geacontrollers.DisplayType.M2M:
				for i := 0; i < vLists.Len(); i++ {
					vLists.Index(i).Interface().(geamodels.M2MModel).LoadM2M()
				}
			case geacontrollers.DisplayType.O2O, geacontrollers.DisplayType.ForeignKey:
				if _, ok := fksList[item.Model]; !ok {
					fksList[item.Model] = reflect.ValueOf(
						c.GEADataBaseQueryList(item.Model, 0, 0, nil, nil, false),
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

func (c *AdaptController) GEADataBaseQueryRow(
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
			case geacontrollers.DisplayType.M2M:
				row.(geamodels.M2MModel).LoadM2M()
			case geacontrollers.DisplayType.O2O, geacontrollers.DisplayType.ForeignKey:
				if _, ok := fksList[item.Model]; !ok {
					fksList[item.Model] = reflect.ValueOf(
						c.GEADataBaseQueryList(item.Model, 0, 0, nil, nil, false),
					).Elem()
				}
				vRow := reflect.ValueOf(row).Elem()
				if vRow.FieldByName(item.Field).IsNil() {
					continue
				}
				vRowIndex := vRow.FieldByName(item.Field).Elem().FieldByName(item.Index).Interface()
				for j := 0; j < fksList[item.Model].Len(); j++ {
					if vRowIndex == fksList[item.Model].
						Index(j).
						Elem().
						FieldByName(item.Index).
						Interface() {
						vRow.FieldByName(item.Field).Set(fksList[item.Model].Index(j))
						break
					}
				}
			}
		}
	}
	return row
}

// TODO M2M O2O FK 处理
func (c *AdaptController) GEADataBaseInsert(model geamodels.Model) (int64, error) {
	ID, err := orm.NewOrm().Insert(model)
	return ID, err
}

func (c *AdaptController) GEADataBaseUpdate(
	model geamodels.Model,
	filters map[string]interface{},
	params map[string]interface{},
) (int64, error) {
	qs := orm.NewOrm().QueryTable(model)
	for k, v := range filters {
		qs = qs.Filter(k, v)
	}
	ID, err := qs.Update((orm.Params)(params))
	return ID, err
}

func (c *AdaptController) GEADataBaseDelete(
	model geamodels.Model, filters map[string]interface{},
) (int64, error) {
	qs := orm.NewOrm().QueryTable(model)
	for k, v := range filters {
		qs = qs.Filter(k, v)
	}
	ID, err := qs.Delete()
	return ID, err
}
