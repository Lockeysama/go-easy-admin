package basecontrollers

import (
	"mime/multipart"
	"net/url"
)

func (c *AdaptController) GetCookie(key string) string {
	return c.Ctx.GetCookie(key)
}

func (c *AdaptController) SetCookie(name string, value string, others ...interface{}) {
	c.Ctx.SetCookie(name, value, others...)
}

func (c *AdaptController) RequestURL() *url.URL {
	return c.Ctx.Request.URL
}

func (c *AdaptController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *AdaptController) RequestHeaderQuery(key string) string {
	return c.Ctx.Request.Header.Get(key)
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
