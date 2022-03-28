package accountmodels

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	basemodels "github.com/lockeysama/go-easy-admin/beego_adapt/models/base"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

func init() {
	orm.RegisterModelWithPrefix("user_", new(User))

	geacontrollers.DefaultValueMakerRegister("passwordMaker", passwordMaker)
}

// Salt 盐
const Salt = "]-O?q{Ht)G2ac7Dj=X'}o@^p.vNusVbn"

// User 用户
type User struct {
	basemodels.NormalModel
	Parent      *User     `orm:"rel(fk);null;description(父账户)"`
	AccountType int8      `orm:"default(0);description(账户类型> 0: root | 1: IAM )"`
	EMail       *string   `orm:"null;size(64);unique;column(email);description(邮箱)" display:"required=false" description:"邮箱"`
	Phone       *string   `orm:"null;size(24);unique;description(电话)" description:"电话"`
	Password    *string   `orm:"null;size(256);description(密码)" description:"密码" gea:"readonly=true;maker=passwordMaker"`
	LastLoginIP *string   `orm:"null;column(last_login_ip);size(32);description(最后登录 IP)" description:"最后登录 IP"`
	LastLogin   time.Time `orm:"auto_now;type(datetime)" display:"dbtype=Datetime" description:"最后登录时间"`
}

func passwordMaker() interface{} {
	return "password"
}
