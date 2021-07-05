package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/client/orm"

	adminmodels "TDCS/models/admin"

	_ "TDCS/models/data"
	_ "TDCS/models/hardware"
	_ "TDCS/models/helper"
	_ "TDCS/models/user"
)

func init() {
	sqlconn, _ := beego.AppConfig.String("sqlconn")
	fmt.Println(sqlconn)
	orm.DefaultTimeLoc = time.UTC
	orm.RegisterDataBase("default", "mysql", sqlconn)

	if v, _ := beego.AppConfig.String("runmode"); v == "dev" {
		orm.Debug = true
	}
	if err := orm.RunSyncdb("default", false, true); err != nil {
		fmt.Println(err)
	}

	adminmodels.RegisterCasbin()
	adminmodels.RegisterRoles()
	adminmodels.AddRolesGroupPolicy()
	adminmodels.CreateSuperUser()
}
