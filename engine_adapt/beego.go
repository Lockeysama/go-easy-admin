package geaadapt

import (
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/lockeysama/go-easy-admin/geadmin"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

// Beego
func InitEngine() {
	geadmin.InitEngine(geadmin.Beego, &AdaptController{})

	GEARouters := geadmin.Routers()

	for routerType, routers := range GEARouters {
		for _, _router := range routers {
			router := _router.([]interface{})
			switch routerType {
			case geadmin.Router:
				if len(router) != 3 {
					panic("router exception")
				}
				controller := router[1].(geacontrollers.Controller)
				controller.SetEngine(&AdaptController{})
				// controller := router[1].(beego.ControllerInterface)
				beego.Router(
					router[0].(string),
					controller.(beego.ControllerInterface),
					router[2].(string),
				)
			case geadmin.AutoRouter:
				group := router[0].(string)
				linkNamespaces := []beego.LinkNamespace{}
				for _, linkNamespace := range router[1].([]interface{}) {
					linkNamespaces = append(
						linkNamespaces,
						beego.NSAutoRouter(linkNamespace.(beego.ControllerInterface)),
					)
				}
				beego.AddNamespace(beego.NewNamespace(
					group,
					linkNamespaces...,
				))
			}
		}
	}
}

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
	geacontrollers.BaseController
}

func (c *AdaptController) Init(ctx *context.Context, controllerName string, actionName string, app interface{}) {
	c.Controller.Init(ctx, controllerName, actionName, app)
}

func (c *AdaptController) SetEngine(controller geacontrollers.Controller) {
	c.BaseController.Controller = controller
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
