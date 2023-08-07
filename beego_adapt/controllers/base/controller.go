package basecontrollers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

type AdaptController struct {
	beego.Controller
	geacontrollers.GEAdminBaseController
}

func (c *AdaptController) Init(
	ctx *context.Context, controllerName string, actionName string, app interface{},
) {
	c.Adapter(app)
	c.Controller.Init(ctx, controllerName, actionName, app)
}

func (c *AdaptController) AccessType() string {
	return c.GEAdminBaseController.AccessType()
}

func (c *AdaptController) Prepare() {
	c.Controller.Prepare()
	c.GEAdminBaseController.Prepare()
}

func (c *AdaptController) Redirect(url string, code int) {
	c.Controller.Redirect(url, code)
}

func (c *AdaptController) SetLayout(layout string) {
	c.Layout = layout
}

func (c *AdaptController) SetTplName(tplName string) {
	c.TplName = tplName
}

func (c *AdaptController) GetController() string {
	controller, _ := c.GetControllerAndAction()
	return controller
}

func (c *AdaptController) ControllerName() string {
	ctrl := c.GetController()
	return strings.ToLower(ctrl[0 : len(ctrl)-10])
}

func (c *AdaptController) GetAction() string {
	_, action := c.Controller.GetControllerAndAction()
	return action
}

func (c *AdaptController) ActionName() string {
	return strings.ToLower(c.GetAction())
}

func (c *AdaptController) SetData(dataType interface{}, data interface{}) {
	c.Data[dataType] = data
}

func (c *AdaptController) GetData() map[interface{}]interface{} {
	return c.Data
}

func (c *AdaptController) ServeJSON(encoding ...bool) {
	c.Controller.ServeJSON(encoding...)
}

func (c *AdaptController) CustomAbort(status int, body string) {
	c.Controller.CustomAbort(status, body)
}

func (c *AdaptController) StopRun() {
	c.Controller.StopRun()
}
