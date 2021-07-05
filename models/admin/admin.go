package adminmodels

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Admin))
}

// Salt 盐
const Salt = "[-O?q{Ht)G2Pf7Dj=X'}o@^p.vNusVbz"

// Admin 管理员
type Admin struct {
	ID          int     `orm:"column(id);auto;pk" description:"ID" json:"id" display:"title=ID;pk=true"`
	UserName    string  `orm:"unique" description:"用户名" json:"username" display:"title=用户名"`
	Password    string  `description:"密码" json:"password" display:"-"`
	RealName    string  `description:"真实姓名" json:"realname" display:"title=真实姓名"`
	Phone       string  `description:"电话" json:"phone" display:"title=电话"`
	Email       string  `description:"电邮" json:"email" display:"title=电邮"`
	Avatar      string  `display:"title=头像;dbtype=File;required=false;meta=admin/avatar/"`
	Status      bool    `description:"状态" json:"status" display:"title=状态"`
	LastLogin   int64   `orm:"auto_now" description:"最后登录时间" json:"lastlogin" display:"title=最后登录时间;dbtype=Datetime"`
	LastIP      string  `description:"最后登录 IP" json:"lastip" display:"title=最后登录 IP"`
	CreatedTime int64   `orm:"auto_now_add" description:"创建时间" json:"createdtime" display:"title=创建时间;dbtype=Datetime"`
	UpdatedTime int64   `description:"修改时间" json:"updatedtime" display:"title=修改时间;dbtype=Datetime"`
	Roles       []*Role `orm:"rel(m2m)" description:"拥有角色" json:"roles" display:"title=拥有角色;showfield=Name"`
}

// CreateSuperUser 创建超管
func CreateSuperUser() {
	o := orm.NewOrm()
	password := generatePasswd() // j+7VPLItBmqe9zo-
	hash := md5.New()
	hash.Write([]byte(password + Salt))
	passwordSaltMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	a := Admin{UserName: "admin", Password: passwordSaltMD5, Status: true}
	if isCreate, _, err := o.ReadOrCreate(&a, "UserName"); err != nil {
		panic("error")
	} else {
		if isCreate {
			role := new([]*Role)
			o.QueryTable(&Role{}).Filter("name", "role_admin").All(role)
			m2m := o.QueryM2M(&a, "Roles")
			m2m.Add(*role)
			fmt.Printf("\n*******************\nCreate Super User:\nUserName: %s\nPassword: %s\n*******************\n", a.UserName, password)
		}
	}
}

// generatePasswd 随机生成密码
func generatePasswd() string {
	const (
		NumStr  = "0123456789"
		CharStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		SpecStr = "<>,.?[]+=-)(*^$#@!"
	)
	var passwd []byte = make([]byte, 16, 16)
	var sourceStr string
	sourceStr = fmt.Sprintf("%s%s%s", NumStr, CharStr, SpecStr)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		index := r.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}
	return string(passwd)
}
