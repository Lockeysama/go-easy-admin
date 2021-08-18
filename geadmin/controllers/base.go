package geacontrollers

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"

	geaconfig "github.com/lockeysama/go-easy-admin/geadmin/config"
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
	"github.com/lockeysama/go-easy-admin/geadmin/utils"
	cache "github.com/lockeysama/go-easy-admin/geadmin/utils/cache"
)

// 消息码
const (
	MSG_OK  = 0
	MSG_ERR = -1
)

type GEAController interface {
	Prepare()

	GetCookie(string) string
	SetCookie(string, string, ...interface{})

	RequestURL() *url.URL
	RequestMethod() string

	RequestQuery(string) string
	RequestParam(string) string
	RequestBody() []byte

	RequestForm() url.Values
	RequestMultipartForm() *multipart.Form

	Redirect(url string, code int)

	ServeJSON(encoding ...bool)
	CustomAbort(status int, body string)

	StopRun()

	SetData(dataType interface{}, data interface{})
	GetData() map[interface{}]interface{}

	SetLayout(layoutName string)
	SetTplName(layoutName string)

	GetController() string
	GetAction() string

	ControllerName() string
	ActionName() string
}

// GEABaseController 控制器基础类
type GEABaseController struct {
	GEAController
	NoAuthAction []string
	User         geamodels.GEAdmin
	APIUser      geamodels.Model
	PageSize     int
	CDNStatic    string
}

// Prepare 前期准备
func (c *GEABaseController) Prepare() {
	c.PageSize = 20
	c.SetData("version", geaconfig.GEAConfig().Version)
	c.SetData("sitename", geaconfig.GEAConfig().SiteName)

	c.SetData("path", c.ControllerName())
	c.SetData("pkField", "")
	c.SetData("CDNStatic", c.CDNStatic)
}

// SideTreeAuth Admin 授权验证
func (c *GEABaseController) SideTreeAuth() {
	sideTree, found := cache.DefaultMemCache().Get(fmt.Sprintf("SideTree%d", c.User.GetID()))
	if found && sideTree != nil { //从缓存取菜单
		sideTree := sideTree.(*[]SideNode)
		c.SetData("SideTree", sideTree)
	} else {
		// 左侧导航栏
		casbinRoles := geamodels.AdminPathPermissions()
		sideTree := SideTree(casbinRoles)
		c.SetData("SideTree", sideTree)
		cache.DefaultMemCache().Set(
			fmt.Sprintf("SideTree%d", c.User.GetID()),
			sideTree,
			cache.DefaultMemCacheExpiration,
		)
	}
}

// Redirect 重定向
func (c *GEABaseController) redirect(url string) {
	c.GEAController.Redirect(url, 302)
	c.StopRun()
}

// Display 加载模板
func (c *GEABaseController) Display(tpl ...string) {
	var name string
	if len(tpl) > 0 {
		name = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		if c.GEAController != nil {
			name = c.ControllerName() + "/" + c.ActionName() + ".html"
		}
	}
	c.SetLayout("public/layout.html")
	c.SetTplName(name)
}

// AjaxMsg ajax返回
func (c *GEABaseController) AjaxMsg(msg interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["message"] = msg
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

// AjaxData ajax返回
func (c *GEABaseController) AjaxData(data interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["data"] = data
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

// AjaxList ajax返回 列表
func (c *GEABaseController) AjaxList(msg interface{}, msgNo int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgNo
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	c.SetData("json", out)
	c.ServeJSON()
	c.StopRun()
}

type filePresignData struct {
	Path string
	URL  string `json:"URL"`
}

type FilePresign struct {
	Code int16
	Data []filePresignData
}

// FilePresign 文件授权
func (c *GEABaseController) FilePresign(method string, paths []string) {
	method = strings.ToUpper(method)
	if method == "GET" {
		filePresign := FilePresign{}
		for _, path := range paths {
			if url, err := utils.PresignRequest(method, path); err != nil {
				c.APIRequestError(400, err.Error())
				return
			} else {
				_filePresignData := filePresignData{}
				_filePresignData.Path = path
				_filePresignData.URL = url
				filePresign.Data = append(filePresign.Data, _filePresignData)
			}
		}
		filePresign.Code = 0
		c.SetData("json", filePresign)
		c.ServeJSON()
	} else if method == "PUT" {
		if url, err := utils.PresignRequest(method, paths[0]); err != nil {
			c.APIRequestError(400, err.Error())
			return
		} else {
			filePresign := FilePresign{}
			filePresign.Code = 0
			filePresign.Data = append(
				filePresign.Data,
				filePresignData{Path: paths[0], URL: url},
			)
			c.SetData("json", filePresign)
			c.ServeJSON()
		}
	} else if method == "POST" {
		if url, err := utils.PresignRequest(method, paths[0]); err != nil {
			c.APIRequestError(400, err.Error())
			return
		} else {
			filePresign := FilePresign{}
			filePresign.Code = 0
			filePresign.Data = append(
				filePresign.Data,
				filePresignData{Path: paths[0], URL: url},
			)
			c.SetData("json", filePresign)
			c.ServeJSON()
		}
	}
}

// APIRequestError API 请求错误
func (c *GEABaseController) APIRequestError(code int, msg ...string) {
	errMsg := ""
	for _, m := range msg {
		errMsg += (m + ". ")
	}
	if errMsg == "" {
		errMsg = "请求错误"
	}
	c.CustomAbort(code, errMsg)
}
