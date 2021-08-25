// @APIVersion 1.0.0
// @Title Flowtime Test API
// @Description Flowtime Test API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import beegoadapt "github.com/lockeysama/go-easy-admin/beego_adapt"

func init() {
	beegoadapt.InjectRouters()
}
