package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro"
	api "github.com/micro/go-micro/api/proto"
	"platform/apis/utils"
	"platform/micro-common/composive"
	"platform/micro-common/log"
	microUtils "platform/micro-common/utils"
	proto "platform/proto/service/data-source"
	"time"
)

var (
	dataSourceConnectError int32 = 9000
	srvName = "go.micro.srv.data-source"
	apiName = "go.micro.api.ds"
	dsCodeParamName = "code"
	dsTypeParamName = "type"
	dataSourceParamName = "data"
	dataSourceRowParamName = "body"
)


type Query struct {
	Client proto.DataSourceService
	logger *log.Logger
}

type DataSourceConfig struct {
	LogPath string
}

type ApiResponse struct {
	Status      		int32
	Message     		string
	Result      		interface{}
	SwiftNumber 		string
	CostTime            int32
}


type LogFormat struct {
	Status              int32
	Path       			string
	RequestTime         string
	CostTime   			int32
	SwiftNumber         string
	ParentSwiftNumber   string
	RawData             string
	FormData            interface{}
	Result   			interface{}
	Message             string
}

func NewDsApiService(cnf *DataSourceConfig) (*Query, error) {
	var err error
	var logger *log.Logger
	if cnf.LogPath == "" {
		logger, err = log.NewLogger()
	} else {
		logger, err = log.NewLogger(cnf.LogPath)
	}
	if err != nil {
		return nil, err
	}
	return &Query{logger: logger}, nil
}

func (d *Query) Data(ctx context.Context, req *api.Request, rsp *api.Response) error {
	startTime := time.Now()
	var apiRuntime LogFormat
	var postParams map[string]string
	var formData map[string]string
	var reqData *proto.QueryDataRequest

	apiRuntime = LogFormat{}
	_, postParams, formData = utils.GetMethodParams(req)
	apiRuntime.FormData = formData
	apiRuntime.RawData = formData[dataSourceRowParamName]
	apiRuntime.SwiftNumber = microUtils.GenSwiftNumber(postParams[dataSourceParamName])
	apiRuntime.ParentSwiftNumber = microUtils.GetParentSwiftNumber(postParams)
	apiRuntime.RequestTime = startTime.Format("2006-01-02 15:04:05")
	apiRuntime.Path = req.Path

	if apiRuntime.RawData == "" {
		reqData = &proto.QueryDataRequest{
			DsType:               postParams[dsTypeParamName],
			DsCode:               postParams[dsCodeParamName],
			Params:               postParams[dataSourceParamName],
			XParentSwiftNumber:   apiRuntime.SwiftNumber,
		}
	} else {
		reqData = &proto.QueryDataRequest{
			DsType:               postParams[dsTypeParamName],
			DsCode:               postParams[dsCodeParamName],
			RawData:              postParams[dataSourceParamName],
			XParentSwiftNumber:   apiRuntime.SwiftNumber,
		}
	}

	response, err := d.Client.QueryData(ctx, reqData)
	apiRuntime.CostTime = int32(time.Now().Sub(startTime).Milliseconds())
	if err != nil {
		apiResponse := ApiResponse{
			Status:      dataSourceConnectError,
			Message:     fmt.Sprintf("data source service connect error: %s", err),
			Result:      "",
			SwiftNumber: apiRuntime.SwiftNumber,
			CostTime:    apiRuntime.CostTime,
		}
		d.writToRsp(rsp, apiRuntime, apiResponse)
	} else {
		apiResponse := ApiResponse{
			Status:      response.Status,
			Message:     response.Message,
			Result:      response.Result,
			SwiftNumber: apiRuntime.SwiftNumber,
			CostTime:    apiRuntime.CostTime,
		}
		d.writToRsp(rsp, apiRuntime, apiResponse)
	}
	return nil
}

func (d *Query) Params(ctx context.Context, req *api.Request, rsp *api.Response) error {
	startTime := time.Now()
	var apiRuntime LogFormat
	var postParams map[string]string
	var allParams map[string]string
	var reqData *proto.QueryParamsRequest

	apiRuntime = LogFormat{}
	_, postParams, allParams = utils.GetMethodParams(req)
	apiRuntime.FormData = allParams
	apiRuntime.RawData = allParams[dataSourceRowParamName]
	apiRuntime.SwiftNumber = microUtils.GenSwiftNumber(postParams[dataSourceParamName])
	apiRuntime.ParentSwiftNumber = microUtils.GetParentSwiftNumber(postParams)
	apiRuntime.RequestTime = startTime.Format("2006-01-02 15:04:05")
	apiRuntime.Path = req.Path

	reqData = &proto.QueryParamsRequest{
		DsType:               postParams[dsTypeParamName],
		DsCode:               postParams[dsCodeParamName],
		XParentSwiftNumber:   apiRuntime.SwiftNumber,
	}
	response, err := d.Client.QueryParams(ctx, reqData)
	apiRuntime.CostTime = int32(time.Now().Sub(startTime).Milliseconds())
	if err != nil {
		apiResponse := ApiResponse{
			Status:      dataSourceConnectError,
			Message:     fmt.Sprintf("data source service connect error: %s", err),
			Result:      "",
			SwiftNumber: apiRuntime.SwiftNumber,
			CostTime:    apiRuntime.CostTime,
		}
		d.writToRsp(rsp, apiRuntime, apiResponse)
		return nil
	} else {
		apiResponse := ApiResponse{
			Status:      response.Status,
			Message:     response.Message,
			Result:      response.Params,
			SwiftNumber: apiRuntime.SwiftNumber,
			CostTime:    apiRuntime.CostTime,
		}
		d.writToRsp(rsp, apiRuntime, apiResponse)
	}
	return nil
}


func (d *Query) writToRsp(rsp *api.Response, apiRuntime LogFormat, apiResponse ApiResponse) {
	apiRspBytes, _ := json.Marshal(apiResponse)
	rsp.Body = string(apiRspBytes)
	apiRuntime.Result = apiResponse.Result
	apiRuntime.Status = apiResponse.Status
	apiRuntime.Message = apiResponse.Message
	_ = d.logger.Log(apiRuntime)
}


func main() {
	cnf := &DataSourceConfig{}
	service, err := composive.GetMicroService(cnf, micro.Name(apiName))
	if err != nil {
		fmt.Printf("create service fail \n%s", err)
		return
	}
	err = composive.InitService(service, cnf)
	if err != nil {
		fmt.Printf("init service fail \n%s", err)
		return
	}
	dsService, err := NewDsApiService(cnf)
	if err != nil {
		fmt.Printf("fail to create data source api service \n%s", err)
		return
	}
	dsService.Client = proto.NewDataSourceService(srvName, service.Client())
	err = service.Server().Handle(
		service.Server().NewHandler(
			dsService,
		),
	)
	if err != nil {
		fmt.Printf("fail to handler api server \n%s", err)
		return
	}
	if err := service.Run(); err != nil {
		fmt.Printf("run api service fail \n%s", err)
		return
	}
}
