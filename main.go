package main

import (
	"PropertyPathPlanning/controllers"
	_ "PropertyPathPlanning/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/",&controllers.MainController{})
	beego.Run()

}

