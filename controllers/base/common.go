package basecontrollers

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"TDCS/utils"

	adminmodels "TDCS/models/admin"
	usermodels "TDCS/models/user"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/client/orm"
	cache "github.com/patrickmn/go-cache"
)

// APIVersion API Allow（约定 API 接口必有版本，版本必须加入 APIVersion，否则将不能正常工作）
var APIVersion = []interface{}{}

// APIAuthFunc API 权限校验函数
var APIAuthFunc func(*BaseController) error

// 消息码
const (
	MSG_OK  = 0
	MSG_ERR = -1
)

// BaseController 控制器基础类
type BaseController struct {
	beego.Controller
	Instance       ControllerRolePolicy
	ControllerName string
	ActionName     string
	NoAuthAction   []string
	User           *adminmodels.Admin
	APIUser        *usermodels.User
	PageSize       int
	CDNStatic      string
}

// APIUserDetail API Get 请求 ID
func (c *BaseController) APIUserDetail(loadRel bool) *usermodels.User {
	user := new(usermodels.User)
	if loadRel {
		orm.NewOrm().QueryTable(c.APIUser).Filter("id", c.APIUser.ID).RelatedSel().One(user)
	} else {
		orm.NewOrm().QueryTable(c.APIUser).Filter("id", c.APIUser.ID).One(user)
	}
	return user
}

// Init 初始化
func (c *BaseController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Controller.Init(ctx, controllerName, actionName, app)
	c.NoAuthAction = *new([]string)
	c.CDNStatic = utils.GetenvFromConfig("cdn_static", "").(string)
}

// Prepare 前期准备
func (c *BaseController) Prepare() {
	c.PageSize = 20
	controllerName, actionName := c.GetControllerAndAction()
	c.ControllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	c.ActionName = strings.ToLower(actionName)
	c.Data["version"], _ = beego.AppConfig.String("version")
	c.Data["siteName"], _ = beego.AppConfig.String("site.name")

	prefix := ""
	if c.Instance != nil {
		prefix = c.Instance.Prefix()
	}
	c.Data["prefix"] = prefix
	c.Data["path"] = c.ControllerName
	c.Data["pkField"] = ""
	c.Data["CDNStatic"] = c.CDNStatic

	if (strings.Compare(c.ControllerName, "apidoc")) != 0 {
		if utils.Contain(strings.Split(c.Ctx.Request.URL.Path[1:], "/")[0], &APIVersion) {
			if APIAuthFunc != nil {
				noAuth := false
				for _, action := range c.NoAuthAction {
					if strings.ToLower(action) == c.ActionName {
						noAuth = true
					}
				}
				if !noAuth {
					if err := APIAuthFunc(c); err != nil {
						actions := []interface{}{"getall", "get", "put", "post", "delete"}
						if utils.Contain(c.ActionName, &actions) {
							c.APIRequestError(401, err.Error())
						} else {
							c.AjaxMsg(err.Error(), MSG_ERR)
						}
					}
				}
			} else {
				panic("APIAuthFunc undefined")
			}
		} else {
			c.auth()
			if c.User != nil {
				c.Data["loginUserName"] = fmt.Sprintf("%s(%s)", c.User.RealName, c.User.UserName)
			}
		}
	}
}

// auth 登录权限验证
func (c *BaseController) auth() {
	arr := strings.Split(c.Ctx.GetCookie("auth"), "|")
	if len(arr) == 2 {
		idStr, password := arr[0], arr[1]
		userID, _ := strconv.Atoi(idStr)
		if userID > 0 {
			var err error

			cheUser, found := utils.DefaultCache().Get("uid" + strconv.Itoa(userID))
			user := &adminmodels.Admin{}
			found = false
			if found && cheUser != nil { //从缓存取用户
				user = cheUser.(*adminmodels.Admin)
			} else {
				o := orm.NewOrm()
				query := o.QueryTable(user)
				filters := map[string]interface{}{"id": userID}
				for key := range filters {
					query = query.Filter(key, filters[key])
				}
				query.One(user)
				o.LoadRelated(user, "Roles")

				if err == nil {
					utils.DefaultCache().Set("uid"+strconv.Itoa(userID), user, cache.DefaultExpiration)
				}
			}
			hash := md5.New()
			hash.Write([]byte(c.GetClientIP() + "|" + user.Password + adminmodels.Salt))
			if err == nil && password == fmt.Sprintf("%x", hash.Sum(nil)) {
				c.User = user
				c.SideTreeAuth()
			}

			//不需要权限检查
			noAuth := "loginin/loginout/getnodes/start/show/ajaxapisave/index/group/public/env/code/apidetail"
			isNoAuth := strings.Contains(noAuth, c.ActionName)

			cr := new([]adminmodels.CasbinRule)
			o := orm.NewOrm()
			roles := []string{}
			for _, r := range user.Roles {
				roles = append(roles, r.Name)
			}
			isHasAuth := false
			prefix := "/"
			if c.Instance != nil {
				prefix = c.Instance.Prefix()
			}
			o.QueryTable(&adminmodels.CasbinRule{}).Filter("V0__in", roles).Filter("V1__contains", fmt.Sprintf("%s/%s/%s", prefix, c.ControllerName, c.ActionName)).All(cr)
			if len(*cr) > 0 {
				isHasAuth = true
			}

			if isHasAuth == false && isNoAuth == false {
				c.AjaxMsg("没有权限", MSG_ERR)
				return
			}
		}
	}

	if (c.User == nil || c.User.ID == 0) && (c.ControllerName != "login" || c.ActionName != "loginin") {
		c.Redirect(beego.URLFor("LoginController.LoginIn"))
	}
}

// SideTreeAuth Admin 授权验证
func (c *BaseController) SideTreeAuth() {
	sideTree, found := utils.DefaultCache().Get("SideTree" + strconv.Itoa(c.User.ID))
	found = false
	if found && sideTree != nil { //从缓存取菜单
		sideTree := sideTree.(*[]*SideNode)
		c.Data["SideTree"] = sideTree
	} else {
		// 左侧导航栏
		casbinRoles := adminmodels.AdminPathPermissions()
		sideTree := SideTree(casbinRoles)
		c.Data["SideTree"] = sideTree
		utils.DefaultCache().Set("SideTree"+strconv.Itoa(c.User.ID), sideTree, cache.DefaultExpiration)
	}
}

// GetClientIP 获取用户 IP 地址
func (c *BaseController) GetClientIP() string {
	s := c.Ctx.Request.RemoteAddr
	l := strings.LastIndex(s, ":")
	return s[0:l]
}

// Redirect 重定向
func (c *BaseController) Redirect(url string) {
	c.Controller.Redirect(url, 302)
	c.StopRun()
}

// Display 加载模板
func (c *BaseController) Display(tpl ...string) {
	var name string
	if len(tpl) > 0 {
		name = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		name = c.ControllerName + "/" + c.ActionName + ".html"
	}
	c.Layout = "public/layout.html"
	c.TplName = name
}

// AjaxMsg ajax返回
func (c *BaseController) AjaxMsg(msg interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["message"] = msg
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

// AjaxData ajax返回
func (c *BaseController) AjaxData(data interface{}, msgNo int) {
	out := make(map[string]interface{})
	out["status"] = msgNo
	out["data"] = data
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

// AjaxList ajax返回 列表
func (c *BaseController) AjaxList(msg interface{}, msgNo int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgNo
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	c.Data["json"] = out
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
func (c *BaseController) FilePresign(method string, paths []string) {
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
		c.Data["json"] = filePresign
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
			c.Data["json"] = filePresign
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
			c.Data["json"] = filePresign
			c.ServeJSON()
		}
	}
}
