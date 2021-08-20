package geacontrollers

type GEAController interface {
	GEARequest

	Prepare()

	Redirect(url string, code int)

	SetLayout(layoutName string)
	SetTplName(layoutName string)

	GetController() string
	ControllerName() string

	GetAction() string
	ActionName() string

	SetData(dataType interface{}, data interface{})
	GetData() map[interface{}]interface{}

	ServeJSON(encoding ...bool)
	CustomAbort(status int, body string)

	StopRun()
}
