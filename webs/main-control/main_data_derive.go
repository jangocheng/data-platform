package main

import (
	"fmt"
	"github.com/micro/go-micro/web"
	"platform/micro-common/composive"
	dataDerive "platform/webs/handler/data-derive"
)

func main() {
	cnfInstance := &dataDerive.DataDeriveConfig{}
	webService, err := composive.GetMicroWebService(
		cnfInstance,
		web.Name("go.micro.web.data-derive"),
	)
	if err != nil {
		fmt.Printf("create web service fail \n%s", err)
		return
	}
	err = composive.InitWebService(webService, cnfInstance)
	if err != nil {
		fmt.Printf("init web service fail \n%s", err)
		return
	}
	webDataDerive, err  := dataDerive.NewDataDeriveWeb(*cnfInstance)
	if err != nil {
		fmt.Printf("fail to create data derive web \n%s", err)
		return
	}
	webService.Handle("/", webDataDerive.Router)
	if err := webService.Run(); err != nil {
		fmt.Printf("fail to run data derive web service \n%s", err)
		return
	}
}