package basecontrollers

import (
	"mime/multipart"
	"net/url"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

type BeegoContext context.Context

func (ctx *BeegoContext) InputQuery(key string) string {
	return ctx.Input.Query(key)
}

func (ctx *BeegoContext) InputParam(key string) string {
	return ctx.Input.Param(key)
}

func (ctx *BeegoContext) InputRequestBody() []byte {
	return ctx.Input.RequestBody
}

func (ctx *BeegoContext) APIVersion() string {
	return strings.Split(ctx.Request.URL.Path[1:], "/")[0]
}

func (ctx *BeegoContext) RequestMethod() string {
	return ctx.Request.Method
}

func (ctx *BeegoContext) RequestForm() url.Values {
	return ctx.Request.Form
}

func (ctx *BeegoContext) RequestMultipartForm() *multipart.Form {
	return ctx.Request.MultipartForm
}

func (ctx *BeegoContext) RequestURL() *url.URL {
	return ctx.Request.URL
}

func (ctx *BeegoContext) RequestRemoteAddr() string {
	return ctx.Request.RemoteAddr
}

type AdaptController struct {
	beego.Controller
	geacontrollers.GEABaseController
}

func (c *AdaptController) Init(ctx *context.Context, controllerName string, actionName string, app interface{}) {
	c.GEABaseController.GEAController = c
	c.Controller.Init(ctx, controllerName, actionName, app)
}

func (c *AdaptController) SetEngine(controller geacontrollers.GEAController) {
	c.GEABaseController.GEAController = controller
}

func (c *AdaptController) Ctx() geacontrollers.Context {
	return (*BeegoContext)(c.Controller.Ctx)
}

func (c *AdaptController) GetCookie(key string) string {
	return c.Controller.Ctx.GetCookie(key)
}

func (c *AdaptController) SetCookie(name string, value string, others ...interface{}) {
	c.Controller.Ctx.SetCookie(name, value, others...)
}

func (c *AdaptController) GetController() string {
	controller, _ := c.Controller.GetControllerAndAction()
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
	c.Controller.Ctx.WriteString(content)
}

func (c *AdaptController) StopRun() {
	c.Controller.StopRun()
}

func (c *AdaptController) Prepare() {
	c.Controller.Prepare()
}

func (c *AdaptController) Get() {
	c.Controller.Get()
}

func (c *AdaptController) Post() {
	c.Controller.Post()
}

func (c *AdaptController) Delete() {
	c.Controller.Delete()
}

func (c *AdaptController) Put() {
	c.Controller.Put()
}

func (c *AdaptController) Head() {
	c.Controller.Head()
}

func (c *AdaptController) Patch() {
	c.Controller.Patch()
}

func (c *AdaptController) Options() {
	c.Controller.Options()
}

func (c *AdaptController) Trace() {
	c.Controller.Trace()
}

func (c *AdaptController) Finish() {
	c.Controller.Finish()
}

func (c *AdaptController) Render() error {
	return c.Controller.Render()
}

func (c *AdaptController) XSRFToken() string {
	return c.Controller.XSRFToken()
}

func (c *AdaptController) CheckXSRFCookie() bool {
	return c.Controller.Ctx.CheckXSRFCookie()
}

func (c *AdaptController) HandlerFunc(fnname string) bool {
	return c.Controller.HandlerFunc(fnname)
}

func (c *AdaptController) URLMapping() {
	c.Controller.URLMapping()
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
