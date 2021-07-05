package adminmodels

import (
	"runtime"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

// CasbinRule Casbin 规则
type CasbinRule struct {
	ID    int    `orm:"column(id)" display:"title=ID"` // 自增主键
	PType string `display:"title=类型"`                  // Policy Type - 用于区分 policy和 group(role)
	V0    string `display:"title=角色"`                  // subject
	V1    string `display:"title=资源"`                  // object
	V2    string `display:"title=权限"`                  // action
	V3    string `display:"title=预留字段"`                // 这个和下面的字段无用，仅预留位置，如果你的不是 \
	V4    string `display:"title=预留字段"`                // 	sub, obj, act的话才会用到
	V5    string `display:"title=预留字段"`                // 	如 sub, obj, act, suf就会用到 V3
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(CasbinRule))
}

// Enforcer Casbin 执行器
var Enforcer *casbin.Enforcer

// Adapter Casbin 适配器
type Adapter struct {
	dataSourceAlias string
	o               orm.Ormer
}

// RegisterCasbin 注册 Casbin 规则
func RegisterCasbin() {
	a := &Adapter{}
	a.o = orm.NewOrm()
	// ❤ + 1s
	runtime.SetFinalizer(a, finalizer)
	// Initialize the model from Go code.
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch(r.obj == p.obj) && regexMatch)r.act == p.act)")

	Enforcer, _ = casbin.NewEnforcer(m, a)
	err := Enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
	Enforcer.EnableAutoSave(true)
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
}

// close 解引用
func (a Adapter) close() {
	a.o = nil
}

// createTable 建表
func (a Adapter) createTable() error {
	return orm.RunSyncdb("default", false, true)
}

// dropTable 删表
func (a Adapter) dropTable() error {
	return orm.RunSyncdb("default", true, true)
}

// loadPolicyLine 载入策略
func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	_, err := a.o.QueryTable("admin_casbin_rule").All(&lines)
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

// savePolicyLine 保存策略
func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a Adapter) SavePolicy(model model.Model) error {
	err := orm.RunSyncdb("default", true, true) // a.dropTable()
	if err != nil {
		return err
	}
	// a = persist.Adapter
	err = orm.RunSyncdb("default", false, true) // a.createTable()
	if err != nil {
		return err
	}

	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	_, err = a.o.InsertMulti(len(lines), lines)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Insert(&line)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Delete(&line, "p_type", "v0", "v1", "v2", "v3", "v4", "v5")
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.PType = ptype
	filter := []string{}
	filter = append(filter, "p_type")
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
		if line.V0 != "" {
			filter = append(filter, "v0")
		}
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
		if line.V1 != "" {
			filter = append(filter, "v1")
		}
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
		if line.V2 != "" {
			filter = append(filter, "v2")
		}
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
		if line.V3 != "" {
			filter = append(filter, "v3")
		}
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
		if line.V4 != "" {
			filter = append(filter, "v4")
		}
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
		if line.V5 != "" {
			filter = append(filter, "v5")
		}
	}

	_, err := a.o.Delete(&line, filter...)
	return err
}

// AdminPathPermissions 获取管理用户列表访问权限
func AdminPathPermissions() map[string][]string {
	path := make(map[string][]string)
	c := new([]*CasbinRule)

	_, err := orm.NewOrm().QueryTable("admin_casbin_rule").Filter("v0", "role_admin").Filter("v1__contains", "/list").OrderBy("id").All(c)
	if err != nil {
		return path
	}

	for _, _c := range *c {
		lv := strings.Split(_c.V1, "/")
		path[lv[1]] = append(path[lv[1]], lv[2])
	}
	return path
}
