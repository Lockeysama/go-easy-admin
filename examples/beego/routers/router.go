// @APIVersion 1.0.0
// @Title Flowtime Test API
// @Description Flowtime Test API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	beegoadapt "github.com/lockeysama/go-easy-admin/beego_adapt"
	accountcontroller "github.com/lockeysama/go-easy-admin/examples/beego/controllers/account"
	applicationcontroller "github.com/lockeysama/go-easy-admin/examples/beego/controllers/application"
)

func init() {
	beegoadapt.InjectRouters()

	beego.AddNamespace(beego.NewNamespace("/account",
		beegoadapt.AutoRegistryRouter(&accountcontroller.AccountController{}),
		beegoadapt.AutoRegistryRouter(&accountcontroller.IAMController{}),
	))

	beego.AddNamespace(beego.NewNamespace("/application",
		beegoadapt.AutoRegistryRouter(&applicationcontroller.ApplicationController{}),
		beegoadapt.AutoRegistryRouter(&applicationcontroller.IoTController{}),
		beegoadapt.AutoRegistryRouter(&applicationcontroller.PlatformController{}),
	))
}
