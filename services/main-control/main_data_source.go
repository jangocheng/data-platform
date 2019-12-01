package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"platform/micro-common/composive"
	dsProto "platform/proto/service/data-source"
	"platform/services/handler/data-source"
)

var srvName = "go.micro.srv.data-source"

func main() {
	cnfInstance := &data_source.DataSourceConfig{}
	service, err := composive.GetMicroService(
		cnfInstance,
		micro.Name(srvName),
		micro.Version("latest"),
		)
	if err != nil {
		fmt.Printf("create service fail \n%s", err)
		return
	}
	err = composive.InitService(service, cnfInstance)
	if err != nil {
		fmt.Printf("init service fail \n%s", err)
		return
	}
	dataSourceService, err := data_source.NewDataSourceService(cnfInstance)
	if err != nil {
		fmt.Printf("create data source service fail \n%s", err)
		return
	}
	defer dataSourceService.Close()
	//go dataSourceService.PP()
	err = dsProto.RegisterDataSourceServiceHandler(service.Server(), dataSourceService)
	if err != nil {
		fmt.Printf("register data source service fail \n%s", err)
		return
	}

	if err := service.Run(); err != nil {
		fmt.Printf("run service fail \n%s", err)
	}
}
