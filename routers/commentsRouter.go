package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:LogFileController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:LogFileController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/data:UsageRecordController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/hardware:HardwareConfigController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/hardware:HardwareConfigController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/helper:IMHelperConfigController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/helper:IMHelperConfigController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:AttributeController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:AttributeController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:AttributeController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:AttributeController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:ConfigController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:ConfigController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:ConfigController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:ConfigController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:DeviceUserController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:DeviceUserController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:DeviceUserController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:DeviceUserController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:SocialController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:SocialController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:UserController"] = append(beego.GlobalControllerRouter["github.com/lockeysama/go-easy-admin/controllers/user:UserController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
