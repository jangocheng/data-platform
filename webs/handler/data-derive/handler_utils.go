package data_derive

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	HttpRequest "platform/common/http"
	"platform/common/utils"
	"strings"
)


func callApiDatas(apiCodes []string, params interface{}, ifCache string, hiddenParam string) (map[string]interface{}, int32, error) {
	var apiResults map[string]interface{}
	var err error
	var response *HttpRequest.Response
	apiResults = make(map[string]interface{}, 0)
	for _, apiCode := range apiCodes {
		var reqCon HttpRequest.RequestConfig
		copyApiNodes := make([]string, 0)
		copy(apiNodes, copyApiNodes)
		nodeSelect := apiNodes[rand.Intn(len(apiNodes))]
		switch params.(type) {
		case string:
			paramsStr := params.(string)
			apiParam := make(map[string]string, 0)
			apiParam["code"] = strings.Trim(apiCode, " ")
			apiParam["data"] = paramsStr
			apiParam["cache"] = ifCache
			apiParam["_platform"] = hiddenParam
			reqCon = HttpRequest.RequestConfig{
				Url: fmt.Sprintf("%s://%s/%s", dataApiScheme, nodeSelect, strings.TrimLeft(dataApiDataPath, "/")),
				Data: apiParam,
				Timeout: 16,
				DisTlsVerify:  true,
			}
		default:
			paramsMap := params.(BeforeWrapperMsg)
			apiParam := make(map[string]string, 0)
			paramsJsonBytes, _ := json.Marshal(paramsMap)
			apiParam["code"] = apiCode
			apiParam["data"] = string(paramsJsonBytes)
			apiParam["cache"] = ifCache
			apiParam["_platform"] = hiddenParam
			reqCon = HttpRequest.RequestConfig{
				Url: fmt.Sprintf("%s://%s/%s", dataApiScheme, nodeSelect, strings.TrimLeft(dataApiDataPath, "/")),
				Data: apiParam,
				Timeout: 16,
				DisTlsVerify:  true,
			}
		}
		request := HttpRequest.NewRequest(reqCon)

		response, err = request.Post()
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, "fail post request")
		}
		rspBody, err := response.Body()
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, "fail get response body")
		}
		apiResponse := &ApiDataResponse{}
		err = json.Unmarshal(rspBody, apiResponse)
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, fmt.Sprintf("fail parse data api: %s response", apiCode))
		}
		if apiResponse.Status != dataApiWebSuccessCode {
			err = errors.New(fmt.Sprintf("data api: %s, return not success code", apiCode))
			return apiResults, apiResponse.Status, err
		}
		apiResults[apiCode] = apiResponse.Result
	}
	return apiResults, SuccessStatus, err
}

func callApiParams(apiCodes []string, hiddenParam string) (map[string]Params, int32, error) {
	var apiResults map[string]Params
	var err error
	var response *HttpRequest.Response
	apiResults = make(map[string]Params, 0)
	for _, apiCode := range apiCodes {
		var reqCon HttpRequest.RequestConfig
		apiNodes := apiNodes
		if len(apiNodes) == 0 {
			return apiResults, CallDataApiError, errors.New("can not find api nodes")
		}
		nodeSelect := apiNodes[rand.Intn(len(apiNodes))]
		apiParam := make(map[string]string, 0)
		apiParam["code"] = strings.Trim(apiCode, " ")
		apiParam["platForm"] = hiddenParam
		reqCon = HttpRequest.RequestConfig{
			Url: fmt.Sprintf("%s://%s/%s", dataApiScheme, nodeSelect, strings.TrimLeft(dataApiParamsPath, "/")),
			Params: apiParam,
			Timeout: 16,
			DisTlsVerify:  true,
		}
		request := HttpRequest.NewRequest(reqCon)

		response, err = request.Get()
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, "fail post request")
		}
		rspBody, err := response.Body()
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, "fail get response body")
		}
		apiResponse := &ApiParamsResponse{}
		err = json.Unmarshal(rspBody, apiResponse)
		if err != nil {
			return apiResults, CallDataApiError, errors.Wrap(err, fmt.Sprintf("fail parse data api: %s response", apiCode))
		}
		if apiResponse.Status != dataApiWebSuccessCode {
			err = errors.New(fmt.Sprintf("data api: %s, return not success code", apiCode))
			return apiResults, apiResponse.Status, err
		}
		apiResults[apiCode] = apiResponse.Result
	}
	return apiResults, SuccessStatus, err
}

func getWrapperList(apiItem *DataDerive) ([]Wrapper, []Wrapper, error) {
	var afterWrapperList []Wrapper
	var logWrapperList []Wrapper

	afterWrapperCodeList := strings.Split(apiItem.AfterWrapperCode, ",")
	logWrapperCodeList := strings.Split(apiItem.LogWrapperCode, ",")

	for _, wrapperCode := range afterWrapperCodeList {
		if wrapperCode == "" {
			continue
		}
		wrapperCode := strings.Trim(wrapperCode, " ")
		wrapper, ok := wrapperMapInstance[wrapperCode]
		if !ok {
			return nil, nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "after" {
			return nil, nil, errors.New(fmt.Sprintf("after wrappers contains not after wrapper: %s", wrapper.WrapperCode))
		}
		afterWrapperList = append(afterWrapperList, wrapper)
	}
	for _, wrapperCode := range logWrapperCodeList {
		if wrapperCode == "" {
			continue
		}
		wrapperCode := strings.Trim(wrapperCode, " ")
		wrapper, ok := wrapperMapInstance[wrapperCode]
		if !ok {
			return nil, nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "log" {
			return nil, nil, errors.New(fmt.Sprintf("log wrappers contains not log wrapper: %s", wrapper.WrapperCode))
		}
		logWrapperList = append(logWrapperList, wrapper)
	}
	return afterWrapperList, logWrapperList, nil
}

func getLogWrapperList(apiItem *DataDeriveSet) ([]Wrapper, error) {
	var logWrapperList []Wrapper

	logWrapperCodeList := strings.Split(apiItem.LogWrapperCode, ",")
	for _, wrapperCode := range logWrapperCodeList {
		if wrapperCode == "" {
			continue
		}
		wrapperCode := strings.Trim(wrapperCode, " ")
		wrapper, ok := wrapperMapInstance[wrapperCode]
		if !ok {
			return nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "log" {
			return nil, errors.New(fmt.Sprintf("log wrappers contains no log wrapper: %s", wrapper.WrapperCode))
		}
		logWrapperList = append(logWrapperList, wrapper)
	}
	return logWrapperList, nil
}

func verifyReqParams(req *http.Request) (bool, error) {
	err := req.ParseForm()
	if err != nil {
		return false, errors.Wrap(err, "invalid params")
	}
	deriveCode := req.Form.Get(deriveCodeParamName)
	if deriveCode == "" {
		return false, ParamsError
	}
	return true, nil
}

func getVerifiedDeriveItem(req *http.Request) (*DataDerive, int32, error) {
	var rspStatus int32
	var rspMessage string
	var deriveItem DataDerive
	var deriveFound bool
	yes, err := verifyReqParams(req)
	if !yes || err == ParamsError {
		rspStatus = ParamsErrorStatus
		if err != nil {
			rspMessage = err.Error()
		}
		rspMessage = "received wrong params"
	} else {
		deriveItem, deriveFound = dataDeriveMapInstance[req.Form.Get(deriveCodeParamName)]
		if !deriveFound {
			rspStatus = NotFundErrorStatus
			rspMessage = fmt.Sprintf("derive code: %s not found", req.Form.Get(deriveCodeParamName))
		} else {
			invalid := deriveItem.Invalid
			if invalid == 1 {
				rspStatus = InvalidStatus
				rspMessage = fmt.Sprintf("derive code: %s invalid", deriveItem.DeriveCode)
			}
		}
	}
	if rspMessage != "" && rspStatus != 0 {
		return nil, rspStatus, errors.New(rspMessage)
	}
	return &deriveItem, 0, nil
}

func getVerifiedDriveSetItem(req *http.Request) (*DataDeriveSet, []*DataDerive, int32, error) {
	var rspStatus int32
	var rspMessage string
	var deriveItems []*DataDerive
	var deriveSetItem DataDeriveSet
	var found bool
	yes, err := verifyReqParams(req)
	if !yes || err == ParamsError {
		rspStatus = ParamsErrorStatus
		if err != nil {
			rspMessage = err.Error()
		}
		rspMessage = "received wrong params"
	} else {
		deriveSetItem, found = dataDeriveSetMapInstance[req.Form.Get(deriveCodeParamName)]
		if !found {
			rspStatus = NotFundErrorStatus
			rspMessage = "deriveSet code can not found"
		}
	}

	if deriveSetItem.DeriveSetCode != ""  {
		deriveCodes := strings.Split(deriveSetItem.DeriveCodes, ",")
		for _, deriveCode := range deriveCodes {
			if strings.Trim(deriveCode, " ") == "" {
				continue
			}
			deriveItem, found := dataDeriveMapInstance[strings.Trim(deriveCode, " ")]
			if !found {
				rspStatus = NotFundErrorStatus
				rspMessage = fmt.Sprintf("derive code: %s not found", deriveCode)
			} else {
				if deriveItem.Invalid == 1 {
					rspStatus = InvalidStatus
					rspMessage = fmt.Sprintf("derive code: %s is invalid", deriveCode)
				} else {
					deriveItems = append(deriveItems, &deriveItem)
				}
			}
		}

	}

	if rspMessage != "" && rspStatus != 0 {
		return nil, nil, rspStatus, errors.New(rspMessage)
	}
	return &deriveSetItem, deriveItems, 0, nil
}

func filterDeriveApiResults(deriveApis string, apiResults map[string]interface{}) map[string]interface{} {
	filResults := make(map[string]interface{})
	deriveApisList := utils.SplitValString(deriveApis)
	for _, deriveApi := range deriveApisList {
		filResults[deriveApi] = apiResults[deriveApi]
	}

	return filResults
}

func getPostMap(req *http.Request) map[string]string {
	var reqParams = make(map[string]string, 0)
	for k, v := range req.Form {
		if len(v) > 0 {
			reqParams[k] = v[0]
		} else {
			reqParams[k] = ""
		}
	}
	return reqParams
}

