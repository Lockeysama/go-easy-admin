package basecontrollers

import (
	"strings"

	"context"

	"github.com/gin-gonic/gin"
	geacontrollers "github.com/lockeysama/go-easy-admin/geadmin/controllers"
)

type RESTful interface {
	Get(*gin.Context)
	Post(*gin.Context)
	Put(*gin.Context)
	Delete(*gin.Context)
}

type IAdaptController interface {
	Init(*context.Context, string, string, interface{})
	geacontrollers.GEAController
	geacontrollers.GEADataBase
	RESTful
}

type AdaptController struct {
	Data    map[interface{}]interface{}
	TplName string
	Ctx     *gin.Context
	IAdaptController
	geacontrollers.GEAdminBaseController
}

func RegisterRouter(c IAdaptController, engine ...*gin.Engine) {
	var _engine *gin.Engine
	if engine == nil {
		_engine = gin.Default()
	}
	c.Init(nil, c.ControllerName(), c.ActionName(), nil)
	group := _engine.Group(c.ControllerName())
	group.GET(c.ControllerName(), c.Get)
	group.POST(c.ControllerName(), c.Get)
	group.PUT(c.ControllerName(), c.Get)
	group.DELETE(c.ControllerName(), c.Get)
}

func (c *AdaptController) Init(
	ctx *context.Context, controllerName string, actionName string, app interface{},
) {
	c.Adapter(app)
}

func (c *AdaptController) Prepare() {
	c.GEAdminBaseController.Prepare()
}

func (c *AdaptController) Redirect(url string, code int) {
	// http.Redirect(c.Writer, req, rURL, code)
}

func (c *AdaptController) SetLayout(layout string) {
	// c.Layout = layout
}

func (c *AdaptController) SetTplName(tplName string) {
	// c.TplName = tplName
}

func (c *AdaptController) GetController() string {
	// controller, _ := c.GetControllerAndAction()
	// return controller
	return ""
}

func (c *AdaptController) ControllerName() string {
	ctrl := c.GetController()
	return strings.ToLower(ctrl[0 : len(ctrl)-10])
}

func (c *AdaptController) GetAction() string {
	// _, action := c.GetControllerAndAction()
	// return action
	return ""
}

func (c *AdaptController) ActionName() string {
	return strings.ToLower(c.GetAction())
}

func (c *AdaptController) SetData(dataType interface{}, data interface{}) {
	// c.Data[dataType] = data
}

func (c *AdaptController) GetData() map[interface{}]interface{} {
	// return c.Data
	return nil
}

func (c *AdaptController) ServeJSON(encoding ...bool) {
	c.ServeJSON(encoding...)
}

func (c *AdaptController) CustomAbort(status int, body string) {
	c.CustomAbort(status, body)
}

func (c *AdaptController) StopRun() {
	c.StopRun()
}
