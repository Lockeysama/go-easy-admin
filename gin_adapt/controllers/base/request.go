package basecontrollers

import (
	"fmt"
	"mime/multipart"
	"net/url"
)

func (c *AdaptController) GetCookie(key string) string {
	cookie, _ := c.Ctx.Cookie(key)
	return cookie
}

func (c *AdaptController) SetCookie(name string, value string, others ...interface{}) {
	var (
		maxAge   int
		path     string
		domain   string
		secure   bool
		httpOnly bool
	)

	c.Ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *AdaptController) RequestURL() *url.URL {
	return c.Ctx.Request.URL
}

func (c *AdaptController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *AdaptController) RequestQuery(key string) string {
	return c.Ctx.Param(key)
}

func (c *AdaptController) RequestParam(key string) string {
	return c.Ctx.Param(key)
}

func (c *AdaptController) RequestBody() []byte {
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

func (c *AdaptController) RequestForm() url.Values {
	return c.Ctx.Request.Form
}

func (c *AdaptController) RequestMultipartForm() *multipart.Form {
	return c.Ctx.Request.MultipartForm
}
