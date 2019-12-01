package data_derive

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	cutils "platform/common/utils"
	"platform/micro-common/utils"
	"strings"
	"time"
)

var (
	ParamsError = errors.New("the params for the derive received are wrong")
)


func (d *DataDeriveWeb) handlerDataGet(rsp http.ResponseWriter, req *http.Request) {

	if strings.Contains(req.RequestURI, "DeriveSet") {
		DeriveSetResult := make(map[string]interface{})

		deriveItemSet, deriveItems, status, err := getVerifiedDriveSetItem(req)
		if err != nil {
			d.setResponse(req, rsp, status, err.Error(), nil)
			return
		}

		logWrapperList, err := getLogWrapperList(deriveItemSet)
		if err != nil {
			d.setResponse(req, rsp, DataDeriveCenterConfigError, err.Error(), logWrapperList)
			return
		}

		apiResults, status, err := d.deriveDataGet(req, deriveItemSet.VerifyParams, deriveItems...)
		if err != nil {
			d.setResponse(req, rsp, status, err.Error(), logWrapperList)
			return
		}

		for _, item := range deriveItems {
			afterWrapperList, _, err := getWrapperList(item)
			if err != nil {
				d.setResponse(req, rsp, DataDeriveCenterConfigError, err.Error(), logWrapperList)
				return
			}

			apiItemResults := filterDeriveApiResults(item.ApiCodes, apiResults)

			result, status, err := d.deriveDataAddress(req, apiItemResults, afterWrapperList, item, deriveItemSet.VerifyParams)
			if err != nil {
				d.setResponse(req, rsp, status, err.Error(), logWrapperList)
				return
			}
			DeriveSetResult[item.DeriveCode] = result
		}
		d.setResponse(req, rsp, status, "success", logWrapperList, DeriveSetResult)
		return

	} else if strings.Contains(req.RequestURI, "Derive") {
		deriveItem, status, err := getVerifiedDeriveItem(req)
		if err != nil {
			d.setResponse(req, rsp, status, err.Error(), nil)
			return
		}

		afterWrapperList, logWrapperList, err := getWrapperList(deriveItem)
		if err != nil {
			d.setResponse(req, rsp, DataDeriveCenterConfigError, err.Error(), logWrapperList)
			return
		}

		dataGet, status, err := d.deriveDataGet(req, deriveItem.VerifyParams, deriveItem)
		if err != nil {
			d.setResponse(req, rsp, status, err.Error(), logWrapperList)
			return
		}

		result, status, err := d.deriveDataAddress(req, dataGet, afterWrapperList, deriveItem, deriveItem.VerifyParams)
		if err != nil {
			d.setResponse(req, rsp, status, err.Error(), logWrapperList)
			return
		}

		d.setResponse(req, rsp, status, "success", logWrapperList, result)
		return
	}
}


func (d *DataDeriveWeb) deriveDataAddress(req *http.Request, apiResults map[string]interface{},
								afterWrapperList []Wrapper, deriveItem *DataDerive, verifyParams string) (interface{}, int32, error) {

	var reqParamsMap = make(map[string]interface{})
	var reqParamsString string
	var afterParams interface{}
	if strings.Contains(req.RequestURI, "verify") {
		reqParamsString = verifyParams
		if reqParamsString == "" {
			return nil, DataDeriveCenterConfigError, errors.New(fmt.Sprintf("verify params not set"))
		}
	} else {
		reqParamsString = req.Form.Get(dataApiParamName)
	}
	errTrans := json.Unmarshal([]byte(reqParamsString), &reqParamsMap)
	if errTrans == nil {
		afterParams = reqParamsMap
	} else {
		afterParams = reqParamsString
	}

	afterWrapperMsg := &AfterWrapperMsg{
		Fail:    false,
		Result:  apiResults,
		Params:  afterParams,
		Item:    deriveItem,
	}
	err := d.JsCustomWrapper(afterWrapperList, afterWrapperMsg, deriveItem.InLineWrapper)
	if err != nil {
		return nil, WrapperExecError, errors.New(fmt.Sprintf("derive: %s execute js code fail, %s",
			deriveItem.DeriveCode, err.Error()))
	}
	if afterWrapperMsg.Fail == true {
		return nil, DataWrapperJudgeError, errors.New(
			fmt.Sprintf("derive: %s wrapper judge fail", deriveItem.DeriveCode))
	}
	return afterWrapperMsg.Result, SuccessStatus, nil
}

func (d *DataDeriveWeb) deriveDataGet(req *http.Request, verifyParams string,
										deriveItems ...*DataDerive) (map[string]interface{}, int32, error) {
	var err error
	var reqParamsString string
	var reqParamsMap BeforeWrapperMsg
	var apiCodes []string
	var apiCodeFilter = make(map[string]interface{}, 0)
	var apiResult = make(map[string]interface{}, 0)

	for _, deriveItem := range deriveItems {
		deriveApiCodes := cutils.SplitValString(deriveItem.ApiCodes)
		if len(deriveApiCodes) == 0 {
			return nil, DataDeriveCenterConfigError,
			errors.New(fmt.Sprintf("derive: %s api codes not defined", deriveItem.DeriveCode))
		}
		for _, deriveApiCode := range deriveApiCodes {
			if _, ok := apiCodeFilter[deriveApiCode]; !ok {
				apiCodes = append(apiCodes, deriveApiCode)
				apiCodeFilter[deriveApiCode] = ""
			}
		}
	}

	reqParamsMap = make(map[string]interface{})
	if strings.Contains(req.RequestURI, "verify") {
		reqParamsString = verifyParams
		if reqParamsString == "" {
			return nil, DataDeriveCenterConfigError, errors.New(fmt.Sprintf("verify params not set"))
		}
	} else {
		reqParamsString = req.Form.Get(dataApiParamName)
	}

	errJson2Map := json.Unmarshal([]byte(reqParamsString), &reqParamsMap)
	var code int32 = 0
	if errJson2Map == nil {
		apiResult, code, err = callApiDatas(apiCodes, reqParamsMap, req.Form.Get(cacheParamName),
			utils.GetPlatformParam(req.Context().Value("swiftNumber").(string)))
	} else {
		apiResult, code, err = callApiDatas(apiCodes, reqParamsString, req.Form.Get(cacheParamName),
			utils.GetPlatformParam(req.Context().Value("swiftNumber").(string)))
	}

	if err != nil {
		return nil, code, err
	}
	return apiResult, code, nil
}


func (d *DataDeriveWeb) handlerDataApiParams(rsp http.ResponseWriter, req *http.Request) {
	deriveItem, status, err := getVerifiedDeriveItem(req)
	if err != nil {
		d.setResponse(req, rsp, status, err.Error(), nil)
		return
	}
	apiCodesStr := deriveItem.ApiCodes
	apiCodesSlice := cutils.SplitValString(apiCodesStr)
	if len(apiCodesSlice) == 0 {
		d.setResponse(req, rsp, DataDeriveCenterConfigError,
			"the api codes for the data derive has not defined", nil)
		return
	}
	var apiParamsResult = make(map[string]Params, 0)
	var code int32 = 0
	apiParamsResult, code, err = callApiParams(apiCodesSlice,
		utils.GetPlatformParam(req.Context().Value("swiftNumber").(string)))
	if err != nil {
		d.setResponse(req, rsp, code, err.Error(), nil)
		return
	}
	d.setResponse(req, rsp, SuccessStatus, "success", nil , apiParamsResult)
	return
}


func (d *DataDeriveWeb) setResponse(req *http.Request, rsp http.ResponseWriter, status int32, msg string,
	logWrapper []Wrapper, result ...interface{}) {
	rspSet := &DeriveResponse{
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
		rspSet.Result = "fail to build response data for data api "
		bodyBytes, _ = json.Marshal(rspSet)
	}
	_, _ = rsp.Write(bodyBytes)
	d.logHandler(req, rspSet, logWrapper)
}


func (d *DataDeriveWeb) logHandler(req *http.Request, rspSet *DeriveResponse, logWrapper []Wrapper) {
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
	})
}




