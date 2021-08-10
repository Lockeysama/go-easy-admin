package geamodels

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Admin), new(AdminRole))
}

// Salt 盐
const Salt = "[-O?q{Ht)G2Pf7Dj=X'}o@^p.vNusVbz"

// Admin 管理员
type Admin struct {
	NormalModel
	UserName  string  `orm:"unique" description:"用户名" display:"title=用户名"`
	Password  string  `description:"密码" display:"-"`
	RealName  string  `description:"真实姓名" display:"title=真实姓名"`
	Phone     string  `description:"电话" display:"title=电话"`
	Email     string  `description:"电邮" display:"title=电邮"`
	Avatar    string  `display:"title=头像;dbtype=File;required=false;meta=admin/avatar/"`
	Status    bool    `description:"状态" display:"title=状态"`
	LastLogin int64   `orm:"auto_now" description:"最后登录时间" display:"title=最后登录时间;dbtype=Datetime"`
	LastIP    string  `orm:"column(last_ip)" description:"最后登录 IP" display:"title=最后登录 IP"`
	Roles     []*Role `orm:"-" description:"拥有角色" display:"title=拥有角色;showfield=Name"`
}

type AdminRole struct {
	NormalModel
	AdminID int64 `orm:"column(admin_id)" description:"Admin ID" json:"admin_id" display:"title=Admin ID"`
	RoleID  int64 `orm:"column(role_id)" description:"Role ID" json:"role_id" display:"title=Role ID"`
}

// CreateSuperUser 创建超管
func CreateSuperUser() {
	o := orm.NewOrm()
	password := generatePasswd() // j+7VPLItBmqe9zo-
	hash := md5.New()
	hash.Write([]byte(password + Salt))
	passwordSaltMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	a := Admin{UserName: "admin", Password: passwordSaltMD5, Status: true}
	if isCreate, adminID, err := o.ReadOrCreate(&a, "UserName"); err != nil {
		panic("error")
	} else {
		if isCreate {
			role := new(Role)
			if err := o.QueryTable(&Role{}).Filter("name", "role_admin").One(role); err != nil {
				panic(err.Error())
			} else {
				adminRole := new(AdminRole)
				adminRole.AdminID = adminID
				adminRole.RoleID = role.ID
				if _, err := o.Insert(adminRole); err != nil {
					panic(err.Error())
				}
				fmt.Printf(""+
					"\n*******************\n"+
					"Create Super User:\n"+
					"UserName: %s\n"+
					"Password: %s\n"+
					"*******************\n",
					a.UserName,
					password,
				)
			}
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
