package main

import (
	// "github.com/casbin/casbin/v2/config"
	_ "github.com/lockeysama/go-easy-admin/models"
	_ "github.com/lockeysama/go-easy-admin/routers"
	"github.com/lockeysama/go-easy-admin/utils"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/lockeysama/go-easy-admin/config"
	"github.com/lockeysama/go-easy-admin/utils/confighelper"
)

func main() {
	config := config.Config{}
	confighelper.LoadConfig(confighelper.ENV, "", &config)
	utils.InitCache()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
