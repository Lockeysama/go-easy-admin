package basecontrollers

import (
	"mime/multipart"
	"net/url"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

type AdaptController struct {
	beego.Controller
	geacontrollers.GEAdminBaseController
}

func (c *AdaptController) Init(ctx *context.Context, controllerName string, actionName string, app interface{}) {
	c.Adapter(app)
	c.Controller.Init(ctx, controllerName, actionName, app)
}

func (c *AdaptController) RequestURL() *url.URL {
	return c.Ctx.Request.URL
}

func (c *AdaptController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *AdaptController) RequestQuery(key string) string {
	return c.Ctx.Input.Query(key)
}

func (c *AdaptController) RequestParam(key string) string {
	return c.Ctx.Input.Param(key)
}

func (c *AdaptController) RequestBody() []byte {
	return c.Ctx.Input.RequestBody
}

func (c *AdaptController) RequestForm() url.Values {
	return c.Ctx.Request.Form
}

func (c *AdaptController) RequestMultipartForm() *multipart.Form {
	return c.Ctx.Request.MultipartForm
}

func (c *AdaptController) GetCookie(key string) string {
	return c.Ctx.GetCookie(key)
}

func (c *AdaptController) SetCookie(name string, value string, others ...interface{}) {
	c.Ctx.SetCookie(name, value, others...)
}

func (c *AdaptController) GetController() string {
	controller, _ := c.GetControllerAndAction()
	return controller
}

func (c *AdaptController) GetAction() string {
	_, action := c.Controller.GetControllerAndAction()
	return action
}

func (c *AdaptController) Redirect(url string, code int) {
	c.Controller.Redirect(url, code)
}

func (c *AdaptController) ServeJSON(encoding ...bool) {
	c.Controller.ServeJSON(encoding...)
}

func (c *AdaptController) CustomAbort(status int, body string) {
	c.Controller.CustomAbort(status, body)
}

func (c *AdaptController) WriteString(content string) {
	c.Ctx.WriteString(content)
}

func (c *AdaptController) StopRun() {
	c.Controller.StopRun()
}

func (c *AdaptController) Prepare() {
	c.Controller.Prepare()
	c.GEAdminBaseController.Prepare()
}

func (c *AdaptController) SetData(dataType interface{}, data interface{}) {
	c.Data[dataType] = data
}

func (c *AdaptController) GetData() map[interface{}]interface{} {
	return c.Data
}

func (c *AdaptController) SetLayout(layout string) {
	c.Layout = layout
}

func (c *AdaptController) SetTplName(tplName string) {
	c.TplName = tplName
}

func (c *AdaptController) ControllerName() string {
	ctrl := c.GetController()
	return strings.ToLower(ctrl[0 : len(ctrl)-10])
}

func (c *AdaptController) ActionName() string {
	return strings.ToLower(c.GetAction())
}

func (c *AdaptController) List() {
	c.GEAdminBaseController.List()
}
