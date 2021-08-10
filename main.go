package main

import (
	adapt "github.com/lockeysama/go-easy-admin/engine_adapt"
	"github.com/lockeysama/go-easy-admin/geadmin/utils"
	"github.com/lockeysama/go-easy-admin/geadmin/utils/cache"
	_ "github.com/lockeysama/go-easy-admin/models"
	_ "github.com/lockeysama/go-easy-admin/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	adapt.InitEngine()
	beego.SetViewsPath("geadmin/views")
	beego.SetStaticPath("/static", "geadmin/static")
	beego.AddFuncMap("FileExt", utils.FileExt)
	// config := config.Config{}
	// confighelper.LoadConfig(confighelper.ENV, "", &config)
	// confighelper.LoadConfig(confighelper.YAML, "./conf/app.yaml", &config)
	cache.InitCache()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
