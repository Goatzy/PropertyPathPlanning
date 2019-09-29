package routers

import (
	"PropertyPathPlanning/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/show", &controllers.MainController{})
	beego.Router("/",&controllers.MainController{})
	beego.Router("/search",&controllers.MainController{})
}
