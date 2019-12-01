package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"platform/common/utils"
	"platform/data-proxy/controllers"
	"platform/data-proxy/model"
)

func main() {
	if len(os.Args) < 2 {
		panic("not receive config path")
	}

	err := utils.LoadJsonConf(os.Args[1], model.GlobalConf)
	if err != nil {
		panic(err.Error())
	}
	model.InitModel()
	if !model.GlobalConf.Debug {
		logs.SetLevel(logs.LevelInfo)
	}
	beego.Router("/ice_wall/*", &controllers.MainController{}, "post:Router")
	beego.Run()
}

