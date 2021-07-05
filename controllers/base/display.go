package basecontrollers

import (
	basemodels "TDCS/models/base"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DisplayItem 管理后台列表、表单显示配置项
type DisplayItem struct {
	Field     string `json:"field"`
	PK        string `json:"pk"`
	Title     string `json:"title"`
	DataType  string `json:"datatype"`
	DBType    string `json:"dbtype"`
	Blur      string `json:"blur"`
	CDN       string `json:"cdn"`
	Required  string `json:"required"`
	Readonly  string `json:"readonly"`
	ShowField string `json:"showfield"`
	Index     string `json:"index"`
	Model     basemodels.Model
	Meta      string      `json:"meta"`
	Help      string      `json:"help"`
	Value     interface{} `json:"value"`
}

var displayItemTagsDefault = map[string]string{
	"field":     "",
	"pk":        "false",
	"title":     "",
	"datatype":  "",
	"dbtype":    "",
	"blur":      "false",
	"cdn":       "false",
	"required":  "true",
	"readonly":  "false",
	"showfield": "ID",
	"index":     "ID",
	"model":     "",
	"meta":      "",
	"help":      "",
	"value":     "",
}

// DisplayItemsCache DisplayItems 缓存
var DisplayItemsCache = make(map[string]*[]DisplayItem)

// DisplayItems 返回管理后台列表表单显示选项
func (c *ManageBaseController) DisplayItems(model basemodels.Model) *[]DisplayItem {
	v := reflect.ValueOf(model).Elem()
	if v.Kind() == reflect.Invalid {
		v = reflect.New(reflect.TypeOf(model).Elem()).Elem()
	}
	t := v.Type()

	var displayItems *[]DisplayItem
	if displayItems, ok := DisplayItemsCache[t.String()]; ok {
		return CopyDisplayItems(displayItems)
	}

	displayItems = new([]DisplayItem)
	for i := 0; i < v.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get("display"), ";")
		if len(tags) < 1 || tags[0] == "-" {
			continue
		}
		tagsMap := make(map[string]string)
		if len(tags) > 0 && tags[0] != "" {
			for _, tag := range tags {
				tagKV := strings.Split(tag, "=")
				tagsMap[tagKV[0]] = tagKV[1]
			}
		}

		for tag, defaultTag := range displayItemTagsDefault {
			if _, ok := tagsMap[tag]; !ok {
				switch tag {
				case "field":
					tagsMap[tag] = t.Field(i).Name
				case "title":
					ormTags := strings.Split(t.Field(i).Tag.Get("orm"), ";")
					for _, ormTag := range ormTags {
						if strings.Contains(ormTag, "description") {
							tagsMap[tag] = ormTag[12 : len(ormTag)-1]
						}
					}
					if _, ok := tagsMap[tag]; !ok {
						tagsMap[tag] = t.Field(i).Name
					}
				case "datatype":
					tagsMap[tag] = t.Field(i).Type.Name()
				case "dbtype":
					switch t.Field(i).Type.Kind() {
					case reflect.String:
						tagsMap[tag] = "Char"
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
						tagsMap[tag] = "Number"
					case reflect.Bool:
						tagsMap[tag] = "Bool"
					case reflect.Slice:
						tagsMap[tag] = "M2M"
					case reflect.Struct:
						if t.Field(i).Type.Name() == "Time" {
							tagsMap[tag] = "Datetime"
						} else {
							tagsMap[tag] = "ForeignKey"
						}
					case reflect.Ptr:
						switch t.Field(i).Type.Elem().Kind() {
						case reflect.String:
							tagsMap[tag] = "Char"
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
							tagsMap[tag] = "Number"
						case reflect.Bool:
							tagsMap[tag] = "Bool"
						case reflect.Slice:
							tagsMap[tag] = "M2M"
						case reflect.Struct:
							if t.Field(i).Type.Name() == "Time" {
								tagsMap[tag] = "Datetime"
							} else {
								tagsMap[tag] = "ForeignKey"
							}
						}
					default:
						panic(fmt.Sprintf("model field \"%s\" type exception", t.Field(i).Name))
					}

				default:
					tagsMap[tag] = defaultTag
				}
			}
		}

		tagsMapJSON, _ := json.Marshal(tagsMap)
		displayItem := new(DisplayItem)
		json.Unmarshal(tagsMapJSON, displayItem)
		if tagsMap["dbtype"] == "M2M" || tagsMap["dbtype"] == "O2O" || tagsMap["dbtype"] == "ForeignKey" {
			switch t.Field(i).Type.Kind() {
			case reflect.Slice:
				st := reflect.New(t.Field(i).Type.Elem())
				displayItem.Model = st.Elem().Interface().(basemodels.Model)
			default:
				st := reflect.New(t.Field(i).Type)
				displayItem.Model = st.Elem().Interface().(basemodels.Model)
			}
		}
		*displayItems = append(*displayItems, *displayItem)
	}
	DisplayItemsCache[t.String()] = CopyDisplayItems(displayItems)

	return displayItems
}

// CopyDisplayItems 深拷贝结构体
func CopyDisplayItems(src *[]DisplayItem) *[]DisplayItem {
	copied := new([]DisplayItem)
	for _, srcItem := range *src {
		newItem := DisplayItem{}
		newItem.Field = srcItem.Field
		newItem.PK = srcItem.PK
		newItem.Title = srcItem.Title
		newItem.DataType = srcItem.DataType
		newItem.DBType = srcItem.DBType
		newItem.Blur = srcItem.Blur
		newItem.CDN = srcItem.CDN
		newItem.Required = srcItem.Required
		newItem.Readonly = srcItem.Readonly
		newItem.ShowField = srcItem.ShowField
		newItem.Index = srcItem.Index
		newItem.Model = srcItem.Model
		newItem.Meta = srcItem.Meta
		newItem.Help = srcItem.Help
		newItem.Value = srcItem.Value
		*copied = append(*copied, newItem)
	}
	return copied
}

// DisplayFields 显示字段，用于查询数据
func DisplayFields(display *[]DisplayItem) []string {
	fields := []string{}
	for _, d := range *display {
		if d.DBType == "M2M" || d.DBType == "O2O" || d.DBType == "ForeignKey" {
			fields = append(fields, fmt.Sprintf("%s__%s", reflect.TypeOf(d.Model).Elem().Name(), d.Index))
		} else {
			fields = append(fields, d.Field)
		}
	}
	return fields
}

// Struct2MapWithHTML 数据结构体转 map
func Struct2MapWithHTML(obj interface{}, display *[]DisplayItem) map[string]interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		for _, item := range *display {
			if item.Field == t.Field(i).Name {
				switch item.DBType {
				case "Bool":
					if v.Field(i).Bool() {
						data[t.Field(i).Name] = "<span style=\"color: #91c799; font-weight: bold;\">True</span>"
					} else {
						data[t.Field(i).Name] = "<span style=\"color: #F581B1; font-weight: bold;\">False</span>"
					}
				case "Datetime":
					if v.Field(i).Type().Name() == "Time" {
						data[t.Field(i).Name] = v.Field(i).Interface()
					} else {
						data[t.Field(i).Name] = time.Unix(v.Field(i).Int(), 0).Format("2006-01-02 15:04:05")
					}
				case "Time":
					if v.Field(i).Type().Name() == "Time" {
						data[t.Field(i).Name] = v.Field(i).Interface()
					} else {
						second := v.Field(i).Int()
						data[t.Field(i).Name] = fmt.Sprintf(
							"%02d:%02d:%02d",
							second/3600,
							(second%3600)/60,
							(second%3600)%60,
						)
					}
				case "ForeignKey", "O2O":
					if v.Field(i).IsNil() {
						continue
					}
					s := reflect.ValueOf(v.Field(i).Interface()).Elem()
					tpl := "<a name=\"%s\" id=\"%s\" index=\"%v\" class=\"layui-btn layui-btn-xs layui-btn-normal layui-btn-radius\" lay-event=\"detail\">%v</a>"
					values := ""
					switch s.FieldByName(item.ShowField).Kind() {
					case reflect.Int, reflect.Int64:
						value := s.FieldByName(item.ShowField).Int()
						switch s.FieldByName(item.Index).Kind() {
						case reflect.Int, reflect.Int64:
							index := s.FieldByName(item.Index).Int()
							values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
						case reflect.String:
							index := s.FieldByName(item.Index).String()
							values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
						}
					case reflect.String:
						value := s.FieldByName(item.ShowField).String()
						switch s.FieldByName(item.Index).Kind() {
						case reflect.Int, reflect.Int64:
							index := s.FieldByName(item.Index).Int()
							values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
						case reflect.String:
							index := s.FieldByName(item.Index).String()
							values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
						}
					}

					data[t.Field(i).Name] = values
				case "M2M":
					if v.Field(i).IsNil() {
						continue
					}
					s := reflect.ValueOf(v.Field(i).Interface())
					tpl := "<a name=\"%s\" id=\"%s\" index=\"%v\" class=\"layui-btn layui-btn-xs layui-btn-normal layui-btn-radius\" lay-event=\"detail\">%v</a>"
					values := ""
					for _i := 0; _i < s.Len(); _i++ {
						switch s.Index(_i).Elem().FieldByName(item.ShowField).Kind() {
						case reflect.Int, reflect.Int64:
							value := s.Index(_i).Elem().FieldByName(item.ShowField).Int()
							switch s.Index(_i).Elem().FieldByName(item.Index).Kind() {
							case reflect.Int, reflect.Int64:
								index := s.Index(_i).Elem().FieldByName(item.Index).Int()
								values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
							case reflect.String:
								index := s.Index(_i).Elem().FieldByName(item.Index).String()
								values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
							}
						case reflect.String:
							value := s.Index(_i).Elem().FieldByName(item.ShowField).String()
							switch s.Index(_i).Elem().FieldByName(item.Index).Kind() {
							case reflect.Int, reflect.Int64:
								index := s.Index(_i).Elem().FieldByName(item.Index).Int()
								values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
							case reflect.String:
								index := s.Index(_i).Elem().FieldByName(item.Index).String()
								values += fmt.Sprintf(tpl, item.Field, item.Index, index, value)
							}
						}
					}

					data[t.Field(i).Name] = values
				default:
					data[t.Field(i).Name] = v.Field(i).Interface()
				}
			}
		}
	}
	return data
}
