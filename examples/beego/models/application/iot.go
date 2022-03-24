package applicationmodels

import (
	"github.com/beego/beego/v2/client/orm"
	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
)

func init() {
	orm.RegisterModelWithPrefix("application_", new(IoT))
}

// IoT 应用 IoT 平台信息
type IoT struct {
	basemodels.NormalModel
	Application   *Application `orm:"rel(fk)"`
	AccessKey     string       `orm:"size(32)"`
	Secret        string       `orm:"size(32)"`
	InstanceID    string       `orm:"size(32)"`
	ProductID     string       `orm:"size(32)"`
	ProductKey    string       `orm:"size(32)"`
	ProductSecret string       `orm:"size(32)"`
	GroupID       string       `orm:"size(32)"`
	ClientID      string       `orm:"size(32)"`
	Host          string       `orm:"size(128)"`
}
