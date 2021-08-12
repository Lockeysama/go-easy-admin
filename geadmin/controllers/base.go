package geacontrollers

import (
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"

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
	Instance     ControllerRolePolicy
	NoAuthAction []string
	User         *geamodels.Admin
	APIUser      geamodels.Model
	PageSize     int
	CDNStatic    string
}

// Prepare 前期准备
func (c *GEABaseController) Prepare() {
	c.PageSize = 20
	c.SetData("version", geaconfig.GEAConfig().Version)
	c.SetData("sitename", geaconfig.GEAConfig().SiteName)

	prefix := ""
	if c.Instance != nil {
		prefix = c.Instance.Prefix()
	}

	c.SetData("prefix", prefix)
	c.SetData("path", c.ControllerName())
	c.SetData("pkField", "")
	c.SetData("CDNStatic", c.CDNStatic)

	if (strings.Compare(c.ControllerName(), "apidoc")) != 0 {
		// if utils.Contain(c.Ctx().APIVersion(), &APIVersion) {
		// 	if APIAuthFunc != nil {
		// 		noAuth := false
		// 		for _, action := range c.NoAuthAction {
		// 			if strings.ToLower(action) == c.ActionName() {
		// 				noAuth = true
		// 			}
		// 		}
		// 		if !noAuth {
		// 			if err := APIAuthFunc(c); err != nil {
		// 				actions := []interface{}{"getall", "get", "put", "post", "delete"}
		// 				if utils.Contain(c.ActionName, &actions) {
		// 					c.APIRequestError(401, err.Error())
		// 				} else {
		// 					c.AjaxMsg(err.Error(), MSG_ERR)
		// 				}
		// 			}
		// 		}
		// 	} else {
		// 		panic("APIAuthFunc undefined")
		// 	}
		// } else {
		c.auth()
		if c.User != nil {
			c.SetData(
				"loginUserName",
				fmt.Sprintf("%s(%s)", c.User.RealName, c.User.UserName),
			)
		}
		// }
	}
}

// auth 登录权限验证
func (c *GEABaseController) auth() {
	arr := strings.Split(c.GetCookie("auth"), "|")
	if len(arr) == 2 {
		idStr, password := arr[0], arr[1]
		userID, _ := strconv.Atoi(idStr)
		if userID > 0 {
			var err error

			cacheUser, found := cache.DefaultMemCache().Get("uid" + strconv.Itoa(userID))
			user := &geamodels.Admin{}
			if found && cacheUser != nil { //从缓存取用户
				user = cacheUser.(*geamodels.Admin)
			} else {
				o := orm.NewOrm()
				query := o.QueryTable(user)
				filters := map[string]interface{}{"id": userID}
				for key := range filters {
					query = query.Filter(key, filters[key])
				}
				if err := query.One(user); err != nil {
					c.AjaxMsg("用户不存在", MSG_ERR)
					return
				}

				adminRoles := new([]*geamodels.AdminRole)
				if _, err := o.QueryTable(&geamodels.AdminRole{}).
					Filter("admin_id", user.ID).
					All(adminRoles); err != nil {
					c.AjaxMsg("查询异常: "+err.Error(), MSG_ERR)
					return
				}
				rolesID := []interface{}{}
				for _, adminRole := range *adminRoles {
					rolesID = append(rolesID, adminRole.RoleID)
				}

				roles := new([]*geamodels.Role)
				if _, err := o.QueryTable(&geamodels.Role{}).
					Filter("id__in", rolesID...).
					All(roles); err != nil {
					c.AjaxMsg("查询异常: "+err.Error(), MSG_ERR)
					return
				} else {
					user.Roles = *roles
					cache.DefaultMemCache().Set(
						"uid"+strconv.Itoa(userID),
						user,
						cache.DefaultMemCacheExpiration,
					)
				}
			}
			hash := md5.New()
			hash.Write([]byte(user.Password + geamodels.Salt))
			if err == nil && password == fmt.Sprintf("%x", hash.Sum(nil)) {
				c.User = user
				c.SideTreeAuth()
			}

			//不需要权限检查
			noAuth := `login/logout/getnodes/start/show/ajaxapisave/index/group/public/env/code/apidetail`
			isNoAuth := strings.Contains(noAuth, c.ActionName())

			cr := new([]geamodels.CasbinRule)
			o := orm.NewOrm()
			roles := []interface{}{}
			for _, r := range user.Roles {
				roles = append(roles, r.Name)
			}
			isHasAuth := false
			prefix := "/"
			if c.Instance != nil {
				prefix = c.Instance.Prefix()
			}
			o.QueryTable(&geamodels.CasbinRule{}).
				Filter("V0__in", roles...).
				Filter(
					"V1__contains",
					fmt.Sprintf("%s/%s/%s", prefix, c.ControllerName(), c.ActionName()),
				).
				All(cr)
			if len(*cr) > 0 {
				isHasAuth = true
			}

			if !isHasAuth && !isNoAuth {
				c.AjaxMsg("没有权限", MSG_ERR)
				return
			}
		}
	}

	if (c.User == nil || c.User.ID == 0) && (c.ControllerName() != "login" || c.ActionName() != "login") {
		c.redirect("/login")
	}
}

// SideTreeAuth Admin 授权验证
func (c *GEABaseController) SideTreeAuth() {
	sideTree, found := cache.DefaultMemCache().Get(fmt.Sprintf("SideTree%d", c.User.ID))
	if found && sideTree != nil { //从缓存取菜单
		sideTree := sideTree.(*[]SideNode)
		c.SetData("SideTree", sideTree)
	} else {
		// 左侧导航栏
		casbinRoles := geamodels.AdminPathPermissions()
		sideTree := SideTree(casbinRoles)
		c.SetData("SideTree", sideTree)
		cache.DefaultMemCache().Set(
			fmt.Sprintf("SideTree%d", c.User.ID),
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
