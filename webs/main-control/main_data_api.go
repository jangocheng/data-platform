package main

import (
	"fmt"
	"github.com/micro/go-micro/web"
	"platform/micro-common/composive"
	dataApi "platform/webs/handler/data-api"
)

func main() {
	cnfInstance := &dataApi.DataApiConfig{}
	webService, err := composive.GetMicroWebService(
		cnfInstance,
		web.Name("go.micro.web.data-api"),
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
	webDataApi, err  := dataApi.NewDataApiWeb(*cnfInstance)
	if err != nil {
		fmt.Printf("fail to create data api web \n%s", err)
		return
	}
	webService.Handle("/", webDataApi.Router)
	if err := webService.Run(); err != nil {
		fmt.Printf("fail to run data api web service \n%s", err)
		return
	}
}