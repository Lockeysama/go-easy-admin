package accountmodels

import (
	"github.com/beego/beego/v2/client/orm"
	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

func init() {
	orm.RegisterModelWithPrefix("user_", new(IAM))

	geacontrollers.DefaultValueMakerRegister("accessKeyMaker", accessKeyMaker)
	geacontrollers.DefaultValueMakerRegister("accessSecretMaker", accessSecretMaker)
}

// IAM 用户 IAM
type IAM struct {
	basemodels.NormalModel
	User      *User  `orm:"rel(fk);description(拥有者)"`
	AccessKey string `orm:"size(32)" gea:"maker=accessKeyMaker;readonly=true"`
	Secret    string `orm:"size(32)" gea:"maker=accessSecretMaker;readonly=true"`
}

func accessKeyMaker() interface{} {
	return "access key"
}

func accessSecretMaker() interface{} {
	return "access secret"
}
