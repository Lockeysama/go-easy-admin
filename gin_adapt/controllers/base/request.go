package basecontrollers

import (
	"fmt"
	"mime/multipart"
	"net/url"
)

func (c *AdaptAdminController) GetCookie(key string) string {
	cookie, _ := c.Ctx.Cookie(key)
	return cookie
}

func (c *AdaptAdminController) SetCookie(name string, value string, others ...interface{}) {
	var (
		maxAge   int
		path     string
		domain   string
		secure   bool
		httpOnly bool
	)

	c.Ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *AdaptAdminController) RequestURL() *url.URL {
	return c.Ctx.Request.URL
}

func (c *AdaptAdminController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *AdaptAdminController) RequestQuery(key string) string {
	return c.Ctx.Param(key)
}

func (c *AdaptAdminController) RequestParam(key string) string {
	return c.Ctx.Param(key)
}

func (c *AdaptAdminController) RequestBody() []byte {
	red, err := c.Ctx.Request.GetBody()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	buf := []byte{}
	_, err = red.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return buf
}

func (c *AdaptAdminController) RequestForm() url.Values {
	return c.Ctx.Request.Form
}

func (c *AdaptAdminController) RequestMultipartForm() *multipart.Form {
	return c.Ctx.Request.MultipartForm
}
