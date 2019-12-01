package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"platform/data-proxy/model"
)



func (c *MainController) Platform() (result []byte, code int32, err error) {
	ProxyResult := model.Response{}
	reqForm := getFormMap(c.Ctx.Request.Form)

	jsonData := make(map[string]interface{})
	_ = json.Unmarshal([]byte(c.GetString("jsonData")), &jsonData)

	if jsonData["api"] != nil {
		var dataParam = make(map[string]interface{})
		var postForm = make(map[string]string)
		var code string
		var dataParamByte []byte
		var platformByte []byte
		var ok bool

		if code, ok = jsonData["api"].(string); !ok {
			return nil, model.ParamsErrorCode, errors.New("api not string")
		}

		for k, v := range jsonData {
			if k == "_platform" || k == "api" {
				continue
			}
			dataParam[k] = v
		}

		dataParamByte, _ = json.Marshal(dataParam)
		platformByte, _ = json.Marshal(jsonData["_platform"])

		postForm["code"] = code
		postForm["data"] = string(dataParamByte)
		postForm["cache"] = reqForm["cache"]
		postForm["_platform"] = string(platformByte)

		url := fmt.Sprintf("%s%s", model.GlobalConf.PlatformHost,
			model.GlobalConf.PlatformApiPath)

		result, err = postData(url, postForm, false)
	} else if jsonData["varset"] != nil || jsonData["var"] != nil {
		var dataParam = make(map[string]interface{})
		var postForm = make(map[string]string)
		var code string
		var dataParamByte []byte
		var platformByte []byte
		var source = make(map[string]interface{})
		var params = make(map[string]interface{})
		var ok bool

		if code, ok = jsonData["varset"].(string); !ok {
			if code, ok = jsonData["var"].(string); !ok {
				return nil, model.ParamsErrorCode, errors.New("varset and var not string")
			}
		}

		for k, v := range jsonData {
			if k == "_platform" || k == "varset" || k == "source" || k == "var"{
				continue
			}
			dataParam[k] = v
		}

		if source, ok = jsonData["source"].(map[string]interface{}); !ok {
			return nil, model.ParamsErrorCode, errors.New("source wrong format")
		}

		if params, ok = source["params"].(map[string]interface{}); !ok {
			return nil, model.ParamsErrorCode, errors.New("params wrong format")
		}

		if times, ok := jsonData["times"]; ok {
			dataParam["time"] = times
		}

		if details, ok := jsonData["details"]; ok {
			dataParam["details"] = details
		}

		for k, v := range params {
			dataParam[k] = v
		}

		dataParamByte, _ = json.Marshal(dataParam)
		platformByte, _ = json.Marshal(jsonData["_platform"])

		postForm["code"] = code
		postForm["data"] = string(dataParamByte)
		postForm["cache"] = reqForm["cache"]
		postForm["_platform"] = string(platformByte)

		url := ""
		if jsonData["varset"] != nil {
			url = fmt.Sprintf("%s%s", model.GlobalConf.PlatformHost,
				model.GlobalConf.PlatformDeriveSetPath)
		} else {
			url = fmt.Sprintf("%s%s", model.GlobalConf.PlatformHost,
				model.GlobalConf.PlatformDerivePath)
		}

		result, err = postData(url, postForm, false)
	}
	if err != nil {
		return nil, model.QueryErrorCode, err
	}
	resultRsp := model.DataPlatformResponse{}
	_ = json.Unmarshal(result,resultRsp)
	if resultRsp.Status == 50000 || resultRsp.Status == 50020 || resultRsp.Status == 30000 ||
		resultRsp.Status == 40000 || resultRsp.Status == 10000 || resultRsp.Status == 60000 ||
		resultRsp.Status == 70000  {
		ProxyResult.Code = model.QueryErrorCode
	} else if resultRsp.Status == 8000 || resultRsp.Result == nil {
		ProxyResult.Code = model.NoResultCode
	} else {
		ProxyResult.Code = model.SuccessCode
	}
	ProxyResult.Result = resultRsp.Result
	ProxyResult.Message = resultRsp.Message
	ProxyResult.SwiftNumber = resultRsp.SwiftNumber
	result, _ = json.Marshal(ProxyResult)

	return result, ProxyResult.Code, err
}
