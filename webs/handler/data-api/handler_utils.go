package data_api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"platform/common/utils"
	"strings"
	"time"
)


func getWrapperList(apiItem *DataApi) ([]Wrapper, []Wrapper, []Wrapper, error) {
	var beforeWrapperList []Wrapper
	var afterWrapperList []Wrapper
	var logWrapperList []Wrapper

	beforeWrapperCodeList := strings.Split(apiItem.BeforeWrapperCode, ",")
	afterWrapperCodeList := strings.Split(apiItem.AfterWrapperCode, ",")
	logWrapperCodeList := strings.Split(apiItem.LogWrapperCode, ",")
	for _, wrapperCode := range beforeWrapperCodeList {
		if wrapperCode == "" {
			continue
		}
		wrapperCode := strings.Trim(wrapperCode, " ")
		wrapper, ok := wrapperMapInstance[wrapperCode]
		if !ok {
			//setResponse(rsp, NotFundErrorStatus, fmt.Sprintf("can not find wrapper code: %s", wrapperCode))
			return nil, nil, nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "before" {
			return nil, nil, nil, errors.New(fmt.Sprintf("before wrappers contains not before wrapper: %s", wrapper.WrapperCode))
		}
		beforeWrapperList = append(beforeWrapperList, wrapper)
	}
	for _, wrapperCode := range afterWrapperCodeList {
		if wrapperCode == "" {
			continue
		}
		wrapperCode := strings.Trim(wrapperCode, " ")
		wrapper, ok := wrapperMapInstance[wrapperCode]
		if !ok {
			return nil, nil, nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "after" {
			return nil, nil, nil, errors.New(fmt.Sprintf("after wrappers contains not after wrapper: %s", wrapper.WrapperCode))
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
			return nil, nil, nil, errors.New(fmt.Sprintf("not find wrapper code: %s", wrapperCode))
		}
		if wrapper.Type != "log" {
			return nil, nil, nil, errors.New(fmt.Sprintf("log wrappers contains no log wrapper: %s", wrapper.WrapperCode))
		}

		logWrapperList = append(logWrapperList, wrapper)
	}
	return beforeWrapperList, afterWrapperList, logWrapperList, nil
}

func verifyReqParams(req *http.Request) (bool, error) {
	err := req.ParseForm()
	if err != nil {
		return false, errors.Wrap(err, "invalid params")
	}
	apiCode := req.Form.Get(apiCodeParamName)
	if apiCode == "" {
		return false, ParamsError
	}
	return true, nil
}

func getVerifiedApiItem(req *http.Request) (*DataApi, int32, error) {
	var rspStatus int32
	var rspMessage string
	var apiItem DataApi
	var apiFound bool
	yes, err := verifyReqParams(req)
	if !yes || err == ParamsError {
		rspStatus = ParamsErrorStatus
		if err != nil {
			rspMessage = err.Error()
		}
		rspMessage = "received wrong params"
	} else {
		apiItem, apiFound = dataApiMapInstance[req.Form.Get(apiCodeParamName)]
		if !apiFound {
			rspStatus = NotFundErrorStatus
			rspMessage = fmt.Sprintf("api code: %s not found", req.Form.Get(apiCodeParamName))
		} else {
			invalid := apiItem.Invalid
			if invalid == 1 {
				rspStatus = InvalidStatus
				rspMessage = fmt.Sprintf("api code: %s invalid", apiItem.ApiCode)
			}
		}
	}
	if rspMessage != "" && rspStatus != 0 {
		return nil, rspStatus, errors.New(rspMessage)
	}
	return &apiItem, 0, nil
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

func (d *DataApiWeb) cacheGet(keyElem ...string) (result interface{}, hit bool, err error ) {
	keyAssem := ""
	for _, key := range keyElem {
		keyAssem = keyAssem + "_" + key
	}
	keyDigest := utils.GetMd5(keyAssem)
	hit, err = d.cacheClient.CheckKey(keyDigest)
	if !hit {
		return "", hit, nil
	}
	resultStr, err := d.cacheClient.GetStringValue(keyDigest)
	var resultMap = make(map[string]interface{})
	err = json.Unmarshal([]byte(resultStr), &resultMap)
	return resultMap["result"], hit, err
}

func (d *DataApiWeb) cacheSave(value interface{}, keyElem ...string) (err error ) {
	keyAssem := ""
	for _, key := range keyElem {
		keyAssem = keyAssem + "_" + key
	}
	keyDigest := utils.GetMd5(keyAssem)
	var valueMap = make(map[string]interface{})
	valueMap["result"] = value
	valueStr, err := json.Marshal(valueMap)
	if err != nil {
		return err
	}
	err = d.cacheClient.SetKeyExpire(keyDigest, string(valueStr), time.Second * time.Duration(d.cacheExpire))
	return err
}
