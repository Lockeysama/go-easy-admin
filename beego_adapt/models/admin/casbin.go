package adminmodels

import (
	"github.com/beego/beego/v2/client/orm"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(geamodels.CasbinRule))
}

// CasbinRule Casbin 规则
type CasbinRuleAdapter struct{}

func (adapter *CasbinRuleAdapter) CreateTable() error {
	return orm.RunSyncdb("default", false, true)
}

func (adapter *CasbinRuleAdapter) DropTable() error {
	return orm.RunSyncdb("default", true, true)
}

func (adapter *CasbinRuleAdapter) Query(filters ...map[string]interface{}) *[]*geamodels.CasbinRule {
	query := orm.NewOrm().QueryTable(&geamodels.CasbinRule{})
	for _, filter := range filters {
		for k, v := range filter {
			query = query.Filter(k, v)
		}
	}
	rules := new([]*geamodels.CasbinRule)
	query.All(rules)
	return rules
}

func (adapter *CasbinRuleAdapter) Insert(r ...*geamodels.CasbinRule) (int64, error) {
	var err error
	for _, row := range r {
		_, _, err = orm.NewOrm().ReadOrCreate(row, "PType", "v0", "v1", "v2", "v3", "v4", "v5")
	}
	return int64(len(r)), err
}

func (adapter *CasbinRuleAdapter) Delete(r *geamodels.CasbinRule, filters ...string) (int64, error) {
	return orm.NewOrm().Delete(r, filters...)
}
