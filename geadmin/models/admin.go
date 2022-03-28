package geamodels

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/lockeysama/go-easy-admin/geadmin/utils"
)

// Salt 盐
const Salt = "[-O?q{Ht)G2Pf7Dj=X'}o@^p.vNusVbz"

var DefaultGEAdminUsername = "admin"

type GEAdminAdapter interface {
	NewGEAdmin(username string, password string) GEAdmin

	QueryWithID(ID int64) GEAdmin
	Administrator() GEAdmin
	ReadOrCreate(admin GEAdmin, field string) (isCreate bool, ID int64, err error)
}

var geadminAdapter GEAdminAdapter

func GetGEAdminAdapter() GEAdminAdapter {
	if geadminAdapter == nil {
		panic("geadminAdapter is nil")
	}
	return geadminAdapter
}

func SetGEAdminAdapter(adapter GEAdminAdapter) {
	geadminAdapter = adapter
}

// Admin 管理员
type GEAdmin interface {
	Model
	GetID() int64
	GetUserName() string
	GetPassword() string
	GetRealName() string
	GetAvatar() string
	GetRoles() []GEARole
	SetRoles([]GEARole)
}

// CreateAdministrator 创建超管
func CreateAdministrator() {
	rootUser := utils.GetenvFromConfig("gea.admin", DefaultGEAdminUsername).(string)
	password := utils.GetenvFromConfig("gea.admin_pwd", generatePasswd()).(string)
	hash := md5.New()
	hash.Write([]byte(password + Salt))
	passwordSaltMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	a := GetGEAdminAdapter().NewGEAdmin(rootUser, passwordSaltMD5)
	if isCreate, _, err := GetGEAdminAdapter().ReadOrCreate(a, "UserName"); err != nil {
		panic("error")
	} else {
		if isCreate {
			fmt.Printf(""+
				"\n*******************\n"+
				"Create Administrator:\n"+
				"UserName: %s\n"+
				"Password: %s\n"+
				"*******************\n\n",
				a.GetUserName(),
				password,
			)
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
	var passwd []byte = make([]byte, 16)

	sourceStr := fmt.Sprintf("%s%s%s", NumStr, CharStr, SpecStr)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		index := r.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}
	return string(passwd)
}
