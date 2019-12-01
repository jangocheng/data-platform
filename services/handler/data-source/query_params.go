package data_source

import (
	"context"
	datasource "platform/proto/service/data-source"
	"reflect"
	"time"
)

func (d *DataSourceService) QueryParams(ctx context.Context, req *datasource.QueryParamsRequest, rsp *datasource.QueryParamsResponse) error {
	startTime := time.Now()
	dsType := req.DsType
	if req.DsCode == "" || req.DsType == "" {
		rsp.Status = ParamsErrorStatus
		rsp.Message = "params ds type and ds code can not be none"
		goto Over
	}
	switch dsType {
	case "file":
		var item FileDataSource
		var got bool
		if item, got = fileDataSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		params := make(map[string]string, 0)
		params["file"] = ""
		must := make([]string, 0)
		rsp.Params = &datasource.Params{}
		rsp.Params.Params = params
		rsp.Params.Must = must
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		break
	case "http":
		var item HttpDataSource
		var got bool
		if item, got = httpDataSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		params, must := SplitParamsStr(item.Params)
		params["_headers"] = ""
		params["_cookies"] = ""
		rsp.Params = &datasource.Params{}
		rsp.Params.Params = params
		rsp.Params.Must = must
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		break
	case "db":
		var item DatabaseSource
		var got bool
		if item, got = databaseSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		params, must := SplitParamsStr(item.Params)
		params["_page"] = ""
		params["_size"] = ""
		params["_orderby"] = ""
		rsp.Params = &datasource.Params{}
		rsp.Params.Params = params
		rsp.Params.Must = must
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		break
	case "kv":
		var item KVDatabaseSource
		var got bool
		if item, got = kvDatabaseSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		params := map[string]string{"key": ""}
		must := []string{"key"}
		rsp.Params = &datasource.Params{}
		rsp.Params.Params = params
		rsp.Params.Must = must
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		break
	case "default":
		rsp.Status = ParamsErrorStatus
		rsp.Message = "ds type not found"
	}
Over:
	costTime := int32(time.Now().Sub(startTime).Milliseconds())
	rsp.CostTime = costTime
	_ = d.logger.Log(logFormat{
		Status:            rsp.Status,
		Path:              reflect.TypeOf(req).Elem().String(),
		CostTime:          costTime,
		RequestTime:       startTime.Format("2006-01-02 15:04:05"),
		ParentSwiftNumber: req.XParentSwiftNumber,
		FormData:          req,
		Result:            rsp.Params,
		Message:           rsp.Message,
	})
	return nil
}



