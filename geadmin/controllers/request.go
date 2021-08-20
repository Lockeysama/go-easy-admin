package geacontrollers

import (
	"mime/multipart"
	"net/url"
)

type GEARequest interface {
	GetCookie(string) string
	SetCookie(string, string, ...interface{})

	RequestURL() *url.URL
	RequestMethod() string

	RequestQuery(string) string
	RequestParam(string) string
	RequestBody() []byte

	RequestForm() url.Values
	RequestMultipartForm() *multipart.Form
}
