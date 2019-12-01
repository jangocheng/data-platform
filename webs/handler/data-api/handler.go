package data_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/pkg/errors"
	"net/http"
	proto "platform/proto/service/data-source"
	"strings"
	"time"
)

var (
	ParamsError = errors.New("invalid params")
)


func (d *DataApiWeb) handlerDataGet(rsp http.ResponseWriter, req *http.Request) {
	var reqParamsMap BeforeWrapperMsg
	var cache Cache
	apiItem, status, err := getVerifiedApiItem(req)
	if err != nil {
		d.setResponse(req, rsp, status, err.Error(), nil, cache)
		return
	}
	var useCache bool
	cacheParam := req.Form.Get(cacheParamName)
	if cacheParam == "1" {
		useCache = true
	}

	beforeWrapperList, afterWrapperList, logWrapperList, err := getWrapperList(apiItem)
	if err != nil {
		d.setResponse(req, rsp, DataApiCenterConfigError, err.Error(), logWrapperList, cache)
		return
	}

	if useCache {
		cache.UseCache = true
		cacheResult, hit, err := d.cacheGet(apiItem.ApiCode, apiItem.BeforeWrapperCode,
			apiItem.AfterWrapperCode, apiItem.DsCode, apiItem.DsType, req.FormValue(dataSourceParamName))
		if err != nil {
			d.setResponse(req, rsp, CacheError, err.Error(), logWrapperList, cache)
			return
		}
		if hit {
			cache.Hit = true
			d.setResponse(req, rsp, SuccessStatus, "success", logWrapperList, cache, cacheResult)
			return
		}
	}

	reqParamsMap = make(map[string]interface{})
	var reqParams string
	if strings.Contains(req.RequestURI, "query") {
		reqParams = req.Form.Get(dataSourceParamName)
	} else if strings.Contains(req.RequestURI, "verify") {
		reqParams = apiItem.VerifyParams
		if reqParams == "" {
			d.setResponse(req, rsp, DataApiCenterConfigError, "verify params not set", logWrapperList, cache)
			return
		}
	} else {
		d.setResponse(req, rsp, ParamsErrorStatus, "url not define", logWrapperList, cache)
		return
	}
	errJson2Map := json.Unmarshal([]byte(reqParams), &reqParamsMap)
	if errJson2Map != nil && beforeWrapperList != nil {
		d.setResponse(req, rsp, DataApiCenterConfigError,
			"not json format string can not set before wrapper", logWrapperList, cache)
		return
	} else {
		err = d.JsWrapper(beforeWrapperList, reqParamsMap)
		if err != nil {
			d.setResponse(req, rsp, WrapperExecError, err.Error(), logWrapperList, cache)
			return
		}
	}
	var rpcReq *proto.QueryDataRequest
	var rpcRsp *proto.QueryDataResponse
	var afterParams interface{}
	var queryResult = make(map[string]interface{})
	dsCodeList := strings.Split(apiItem.DsCode, ",")
	for _, code := range dsCodeList {
		dsCode := strings.Trim(code, " ")
		if errJson2Map == nil {
			afterParams = reqParamsMap
			paramsStr, _ := json.Marshal(reqParamsMap)
			rpcReq = &proto.QueryDataRequest{
				DsType:               apiItem.DsType,
				DsCode:               dsCode,
				Params:               string(paramsStr),
				XParentSwiftNumber:	  req.Context().Value("swiftNumber").(string),
			}
		} else {
			afterParams = req.Form.Get(dataSourceParamName)
			rpcReq = &proto.QueryDataRequest{
				DsType:               apiItem.DsType,
				DsCode:               dsCode,
				RawData:              req.Form.Get(dataSourceParamName),
				XParentSwiftNumber:	  req.Context().Value("swiftNumber").(string),
			}
		}
		rpcRsp, err = d.dsClient.QueryData(context.TODO(), rpcReq, client.WithDialTimeout(time.Second *8))
		if err != nil {
			d.setResponse(req, rsp, DataSourceError, errors.Wrap(err, fmt.Sprintf("query ds data: %s fail", dsCode)).Error(),
				logWrapperList, cache)
			return
		}
		if rpcRsp.Status != dataSourceRpcServiceSuccessCode {
			d.setResponse(req, rsp, DataSourceError, fmt.Sprintf("data source: %s return error code: %d, error message: %s",
				dsCode, rpcRsp.Status, rpcRsp.Message), logWrapperList, cache)
			return
		}
		queryResult[dsCode] = rpcRsp.Result
	}
	var afterResult interface{}
	if len(dsCodeList) <= 1 {
		afterResult =  rpcRsp.Result
	} else {
		afterResult =  queryResult
	}
	afterWrapperMsg := &AfterWrapperMsg{
		Fail:    false,
		Params:  afterParams,
		Result:  afterResult,
		Item:    apiItem,
	}
	err = d.JsCustomWrapper(afterWrapperList, afterWrapperMsg)
	if err != nil {
		d.setResponse(req, rsp, WrapperExecError, err.Error(), logWrapperList, cache)
		return
	}
	if afterWrapperMsg.Fail == true {
		d.setResponse(req, rsp, DataWrapperJudgeError,
			fmt.Sprintf("api: %s wrapper judge fail", apiItem.ApiCode),
			logWrapperList, cache)
		return
	}

	if useCache {
		err := d.cacheSave(afterWrapperMsg.Result, apiItem.ApiCode, apiItem.BeforeWrapperCode,
			apiItem.AfterWrapperCode, apiItem.DsCode, apiItem.DsType, req.FormValue(dataSourceParamName))
		if err != nil {
			d.setResponse(req, rsp, CacheError, err.Error(), logWrapperList, cache)
			return
		}
	}

	d.setResponse(req, rsp, SuccessStatus, "success", logWrapperList, cache, afterWrapperMsg.Result)
	return
}


func (d *DataApiWeb) handlerDataSourceParams(rsp http.ResponseWriter, req *http.Request) {
	apiItem, status, err := getVerifiedApiItem(req)
	if err != nil {
		d.setResponse(req, rsp, status, err.Error(), nil, Cache{})
		return
	}
	var rpcReq *proto.QueryParamsRequest
	var rpcRsp *proto.QueryParamsResponse
	dsCodeList := strings.Split(apiItem.DsCode, ",")
	dsParams := make(map[string]*proto.Params)
	for _, code := range dsCodeList {
		dsCode := strings.Trim(code, " ")
		rpcReq = &proto.QueryParamsRequest{
			DsType:               apiItem.DsType,
			DsCode:               dsCode,
			XParentSwiftNumber:	  req.Context().Value("swiftNumber").(string),
		}
		rpcRsp, err = d.dsClient.QueryParams(context.TODO(), rpcReq, client.WithDialTimeout(time.Second *8))
		if err != nil {
			d.setResponse(req, rsp, DataSourceError, errors.Wrap(err, fmt.Sprintf("query data source: %s params fail", dsCode)).Error(),
				nil, Cache{})
			return
		}
		if rpcRsp.Status != dataSourceRpcServiceSuccessCode {
			d.setResponse(req, rsp, DataSourceError, fmt.Sprintf("data source: %s return error code: %d, error message: %s",
				dsCode, rpcRsp.Status, rpcRsp.Message), nil, Cache{})
			return
		}
		dsParams[dsCode] = rpcRsp.Params
	}

	d.setResponse(req, rsp, SuccessStatus, "success", nil , Cache{}, dsParams)
	return
}


func (d *DataApiWeb) setResponse(req *http.Request, rsp http.ResponseWriter, status int32, msg string,
	logWrapper []Wrapper, cache Cache, result ...interface{}) {
	rspSet := &ApiResponse{
		Status:		status,
		Message: 	msg,
		SwiftNumber: req.Context().Value("swiftNumber").(string),
		CostTime:  int32(time.Now().Sub(req.Context().Value("startTime").(time.Time)).Milliseconds()),
	}
	if len(result) != 0 {
		rspSet.Result = result[0]
	}

	bodyBytes, err := json.Marshal(rspSet)
	if err != nil {
		rspSet.Result = "fail to build response data for data api"
		bodyBytes, _ = json.Marshal(rspSet)
	}
	_, _ = rsp.Write(bodyBytes)
	d.logHandler(req, rspSet, logWrapper, cache)
}

func (d *DataApiWeb) logHandler(req *http.Request, rspSet *ApiResponse, logWrapper []Wrapper, cache Cache) {
	reqParams := getPostMap(req)
	logWrapperMsg := &LogWrapperMsg{
		Params: reqParams,
		Result: rspSet.Result,
	}
	wrapperLogTime := d.LogWrapper(logWrapper, logWrapperMsg)
	_ = d.logger.Log(LogFormat{
		Status:            rspSet.Status,
		Path:              req.URL.Path,
		RequestTime:       req.Context().Value("startTime").(time.Time).Format("2006-01-02 15:04:05"),
		CostTime:          rspSet.CostTime,
		SwiftNumber:       rspSet.SwiftNumber,
		ParentSwiftNumber: req.Context().Value("parentSwiftNumber").(string),
		RawData:           "",
		FormData:          logWrapperMsg.Params,
		Result:            logWrapperMsg.Result,
		LogWrapperTime:    wrapperLogTime,
		Message:           rspSet.Message,
		Cache:             cache,
	})
}





