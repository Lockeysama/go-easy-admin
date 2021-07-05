package basecontrollers

import (
	"encoding/json"
	"reflect"
)

// List 管理后台列表模板渲染
func (c *ManageBaseController) List() {
	c.Data["pageTitle"] = "列表"
	items := c.DisplayItems(c.Model)
	c.Data["display"] = items
	c.makeListPK(items)
	c.Layout = "public/layout.html"
	c.TplName = "public/list.html"
}

// Table 获取管理后台列表数据
func (c *ManageBaseController) Table() {
	// 列表
	page, err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 30
	}

	var lists interface{}

	orderByStr := c.Ctx.Input.Query("order_by")
	listOrderBy := make(map[string]string)
	if orderByStr != "" {
		if err := json.Unmarshal([]byte(orderByStr), &listOrderBy); err != nil {
			c.AjaxMsg("order exception", MSG_ERR)
			return
		}
	}

	query := c.Ctx.Input.Query("query")
	listFilter := make(map[string]interface{})

	if query != "" {
		listFilter = c.ListFilter()
		lists = c.QueryList(c.Model, page, limit, listFilter, listOrderBy, true)
	} else {
		lists = c.QueryList(c.Model, page, limit, nil, listOrderBy, true)
	}

	b, _ := json.Marshal(lists)
	var resultMap = new([]map[string]interface{})
	json.Unmarshal(b, resultMap)

	c.PageSize = limit
	listsValue := reflect.ValueOf(lists).Elem()
	list := make([]map[string]interface{}, listsValue.Len())
	for i := 0; i < listsValue.Len(); i++ {
		x := Struct2MapWithHTML(
			listsValue.Index(i).Elem().Interface(), c.DisplayItems(c.Model),
		)
		list[i] = x
	}
	if query != "" {
		if count, err := c.QueryCount(c.Model, listFilter); err != nil {
			c.AjaxMsg("查询失败", MSG_ERR)
		} else {
			c.AjaxList("成功", MSG_OK, count, list)
		}
	} else {
		if count, err := c.QueryCount(c.Model, nil); err != nil {
			c.AjaxMsg("查询失败", MSG_ERR)
		} else {
			c.AjaxList("成功", MSG_OK, count, list)
		}
	}
}

func (c *ManageBaseController) ListFilter() map[string]interface{} {
	queryStr := c.Ctx.Input.Query("query")
	if queryStr == "" {
		c.AjaxMsg("请输入查询条件", MSG_ERR)
		return nil
	}

	query := new([]map[string]interface{})
	if err := json.Unmarshal([]byte(queryStr), query); err != nil {
		c.AjaxMsg("查询条件异常", MSG_ERR)
		return nil
	}

	if len(*query) == 0 {
		c.AjaxMsg("查询条件异常", MSG_ERR)
		return nil
	}

	filters := make(map[string]interface{})

	for _, q := range *query {
		ok := false
		field := ""
		if field, ok = q["field"].(string); !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}
		expression := ""
		if expression, ok = q["exp"].(string); !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}
		var value interface{}
		if value, ok = q["value"]; !ok {
			c.AjaxMsg("查询条件异常", MSG_ERR)
			return nil
		}

		switch expression {
		case "eq":
			filters[field] = value
		case "ne":
			filters[field+"__iexact"] = value
		case "gt":
			filters[field+"__gt"] = value
		case "lt":
			filters[field+"__lt"] = value
		case "gte":
			filters[field+"__gte"] = value
		case "lte":
			filters[field+"__lte"] = value
		case "is_contains":
			filters[field+"__contains"] = value
		case "not_contains":
			filters[field+"__icontains"] = value
		}
	}
	return filters
}
