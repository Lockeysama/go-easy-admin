package adminmodels

import (
	"fmt"

	ginmodels "github.com/lockeysama/go-easy-admin/gin_adapt/models"

	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

func init() {
	if err := ginmodels.DB().AutoMigrate(&CasbinRule{}); err != nil {
		fmt.Println(err.Error())
	}
}

type CasbinRule struct {
	geamodels.CasbinRule `gorm:"embedded"`
}

func (t CasbinRule) TableName() string {
	return "admin_casbin_rule"
}

func NewCasbinRule(r *geamodels.CasbinRule) *CasbinRule {
	cr := new(CasbinRule)
	cr.ID = r.ID
	cr.PType = r.PType
	cr.V0 = r.V0
	cr.V1 = r.V1
	cr.V2 = r.V2
	cr.V3 = r.V3
	cr.V4 = r.V4
	cr.V5 = r.V5
	return cr
}

// CasbinRule Casbin 规则
type CasbinRuleAdapter struct{}

func (adapter *CasbinRuleAdapter) CreateTable() error {
	return nil
}

func (adapter *CasbinRuleAdapter) DropTable() error {
	return nil
}

func (adapter *CasbinRuleAdapter) Query(filters ...map[string]interface{}) *[]*geamodels.CasbinRule {
	rules := new([]*geamodels.CasbinRule)
	ginmodels.DB().Model(&CasbinRule{}).Where(filters).Find(rules)
	return rules
}

func (adapter *CasbinRuleAdapter) Insert(r ...*geamodels.CasbinRule) (int64, error) {
	var err error
	for _, row := range r {
		_row := NewCasbinRule(row)
		result := ginmodels.DB().Create(_row)
		if result.Error != nil {
			return 0, err
		}
	}
	return int64(len(r)), err
}

func (adapter *CasbinRuleAdapter) Delete(r *geamodels.CasbinRule, filters ...string) (int64, error) {
	row := NewCasbinRule(r)
	result := ginmodels.DB().Where(filters).Delete(row)
	return result.RowsAffected, result.Error
}
