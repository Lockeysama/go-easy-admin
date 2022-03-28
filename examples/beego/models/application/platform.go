package applicationmodels

import (
	"github.com/beego/beego/v2/client/orm"
	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
)

func init() {
	orm.RegisterModelWithPrefix("application_", new(Platform))
}

// Platform 应用 Platform 平台信息
type Platform struct {
	basemodels.NormalModel
	Name        string         `orm:"size(32)"`
	Desc        string         `orm:"size(32)"`
	Application []*Application `orm:"reverse(many)"`
}
