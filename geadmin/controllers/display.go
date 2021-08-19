package geacontrollers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
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
	Model     geamodels.Model
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

func fieldParse(field reflect.StructField) (tagsMaps []map[string]string) {
	tagsMap := make(map[string]string)

	tags := strings.Split(field.Tag.Get("display"), ";")
	if len(tags) < 1 || tags[0] == "-" {
		return
	}

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
				tagsMap[tag] = field.Name
			case "title":
				ormTags := strings.Split(field.Tag.Get("orm"), ";")
				for _, ormTag := range ormTags {
					if strings.Contains(ormTag, "description") {
						tagsMap[tag] = ormTag[12 : len(ormTag)-1]
					}
				}
				if _, ok := tagsMap[tag]; !ok {
					tagsMap[tag] = field.Name
				}
			case "datatype":
				tagsMap[tag] = field.Type.Name()
			case "dbtype":
				switch field.Type.Kind() {
				case reflect.String:
					tagsMap[tag] = "Char"
				case reflect.Int,
					reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Float32, reflect.Float64:
					tagsMap[tag] = "Number"
				case reflect.Bool:
					tagsMap[tag] = "Bool"
				case reflect.Slice:
					tagsMap[tag] = "M2M"
				case reflect.Struct:
					if field.Name == field.Type.Name() {
						v := reflect.New(field.Type).Elem()
						t := v.Type()
						for i := 0; i < v.NumField(); i++ {
							subTagsMap := fieldParse(t.Field(i))
							tagsMaps = append(tagsMaps, subTagsMap...)
						}
						return
					} else if field.Type.Name() == "Time" {
						tagsMap[tag] = "Datetime"
					} else {
						tagsMap[tag] = "ForeignKey"
					}
				case reflect.Ptr:
					switch field.Type.Elem().Kind() {
					case reflect.String:
						tagsMap[tag] = "Char"
					case reflect.Int,
						reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Float32, reflect.Float64:
						tagsMap[tag] = "Number"
					case reflect.Bool:
						tagsMap[tag] = "Bool"
					case reflect.Slice:
						tagsMap[tag] = "M2M"
					case reflect.Struct:
						if field.Type.Name() == "Time" {
							tagsMap[tag] = "Datetime"
						} else {
							tagsMap[tag] = "ForeignKey"
						}
					}
				default:
					panic(fmt.Sprintf("model field \"%s\" type exception", field.Name))
				}

			default:
				tagsMap[tag] = defaultTag
			}
		}
	}
	tagsMaps = append(tagsMaps, tagsMap)
	return
}

// Display 加载模板
func (c *GEAdminBaseController) Display(tpl ...string) {
	var name string
	if len(tpl) > 0 {
		name = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		if c.GEAController != nil {
			name = c.ControllerName() + "/" + c.ActionName() + ".html"
		}
	}
	c.SetLayout("public/layout.html")
	c.SetTplName(name)
}

// DisplayItems 返回管理后台列表表单显示选项
func (c *GEAdminBaseController) DisplayItems(model geamodels.Model) *[]DisplayItem {
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
		tagsMaps := fieldParse(t.Field(i))
		for _, tagsMap := range tagsMaps {
			tagsMapJSON, _ := json.Marshal(tagsMap)
			displayItem := new(DisplayItem)
			json.Unmarshal(tagsMapJSON, displayItem)

			if tagsMap["dbtype"] == "M2M" || tagsMap["dbtype"] == "O2O" || tagsMap["dbtype"] == "ForeignKey" {
				switch t.Field(i).Type.Kind() {
				case reflect.Slice:
					st := reflect.New(t.Field(i).Type.Elem())
					displayItem.Model = st.Elem().Interface().(geamodels.Model)
				default:
					st := reflect.New(t.Field(i).Type)
					displayItem.Model = st.Elem().Interface().(geamodels.Model)
				}
			}

			*displayItems = append(*displayItems, *displayItem)
		}
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
func Struct2MapWithHTML(obj *map[string]interface{}, display *[]DisplayItem) map[string]interface{} {
	var data = make(map[string]interface{})
	var ok bool
	for _, item := range *display {
		if item.Value, ok = (*obj)[item.Field]; !ok {
			item.Value = nil
		}
		switch item.DBType {
		case "Bool":
			if item.Value.(bool) {
				data[item.Field] = "<span style=\"color: #91c799; font-weight: bold;\">True</span>"
			} else {
				data[item.Field] = "<span style=\"color: #F581B1; font-weight: bold;\">False</span>"
			}
		case "Datetime":
			switch item.Value.(type) {
			case int64:
				data[item.Field] = time.Unix(item.Value.(int64), 0).Format("2006-01-02 15:04:05")
			case string:
				data[item.Field] = item.Value.(string)
			}
		case "Time":
			second := item.Value.(int)
			data[item.Field] = fmt.Sprintf(
				"%02d:%02d:%02d",
				second/3600,
				(second%3600)/60,
				(second%3600)%60,
			)
		case "ForeignKey", "O2O":
			if item.Value == nil {
				continue
			}
			s := reflect.ValueOf(item.Value).Elem()
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

			data[item.Field] = values
		case "M2M":
			if item.Value == nil {
				continue
			}
			s := reflect.ValueOf(item.Value)
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

			data[item.Field] = values
		default:
			data[item.Field] = item.Value
		}
	}
	return data
}
