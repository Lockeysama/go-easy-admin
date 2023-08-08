package basecontrollers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

type AdaptAdminController struct {
	beego.Controller
	geacontrollers.GEAdminBaseController
}

func (c *AdaptAdminController) Init(
	ctx *context.Context, controllerName string, actionName string, app interface{},
) {
	c.Adapter(app)
	c.Controller.Init(ctx, controllerName, actionName, app)
}

func (c *AdaptAdminController) AccessType() string {
	return geacontrollers.AccessTypeCookie
}

func (c *AdaptAdminController) Prepare() {
	c.Controller.Prepare()
	c.GEAdminBaseController.Prepare()
}

func (c *AdaptAdminController) Redirect(url string, code int) {
	c.Controller.Redirect(url, code)
}

func (c *AdaptAdminController) SetLayout(layout string) {
	c.Layout = layout
}

func (c *AdaptAdminController) SetTplName(tplName string) {
	c.TplName = tplName
}

func (c *AdaptAdminController) GetController() string {
	controller, _ := c.GetControllerAndAction()
	return controller
}

func (c *AdaptAdminController) ControllerName() string {
	ctrl := c.GetController()
	return strings.ToLower(ctrl[0 : len(ctrl)-10])
}

func (c *AdaptAdminController) GetAction() string {
	_, action := c.Controller.GetControllerAndAction()
	return action
}

func (c *AdaptAdminController) ActionName() string {
	return strings.ToLower(c.GetAction())
}

func (c *AdaptAdminController) SetData(dataType interface{}, data interface{}) {
	c.Data[dataType] = data
}

func (c *AdaptAdminController) GetData() map[interface{}]interface{} {
	return c.Data
}

func (c *AdaptAdminController) ServeJSON(encoding ...bool) {
	c.Controller.ServeJSON(encoding...)
}

func (c *AdaptAdminController) CustomAbort(status int, body string) {
	c.Controller.CustomAbort(status, body)
}

func (c *AdaptAdminController) StopRun() {
	c.Controller.StopRun()
}

type AdaptAPIController struct {
	AdaptAdminController
}

func (c *AdaptAPIController) AccessType() string {
	return geacontrollers.AccessTypeJWT
}
