package main

import (
	_ "github.com/lockeysama/go-easy-admin/examples/beego/models"
	_ "github.com/lockeysama/go-easy-admin/examples/beego/routers"
	"github.com/lockeysama/go-easy-admin/geadmin/utils"
	"github.com/lockeysama/go-easy-admin/geadmin/utils/cache"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
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
