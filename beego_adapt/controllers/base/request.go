package basecontrollers

import (
	"mime/multipart"
	"net/url"
)

func (c *AdaptAdminController) GetCookie(key string) string {
	return c.Ctx.GetCookie(key)
}

func (c *AdaptAdminController) SetCookie(name string, value string, others ...interface{}) {
	c.Ctx.SetCookie(name, value, others...)
}

func (c *AdaptAdminController) RequestURL() *url.URL {
	return c.Ctx.Request.URL
}

func (c *AdaptAdminController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *AdaptAdminController) RequestHeaderQuery(key string) string {
	return c.Ctx.Request.Header.Get(key)
}

func (c *AdaptAdminController) RequestQuery(key string) string {
	return c.Ctx.Input.Query(key)
}

func (c *AdaptAdminController) RequestParam(key string) string {
	return c.Ctx.Input.Param(key)
}

func (c *AdaptAdminController) RequestBody() []byte {
	return c.Ctx.Input.RequestBody
}

func (c *AdaptAdminController) RequestForm() url.Values {
	return c.Ctx.Request.Form
}

func (c *AdaptAdminController) RequestMultipartForm() *multipart.Form {
	return c.Ctx.Request.MultipartForm
}
