package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"github.com/satori/go.uuid"
	"platform/common/utils"
	"platform/data-proxy/model"
	"strings"
)

type StrategyListResponse struct {
	Code       string  		`json:"code"`
	Msg        string   	`json:"msg"`
	Data       interface{}  `json:"data"`
}

type StrategyFlowResponse struct {
	Code       string  		`json:"code"`
	Msg        string   	`json:"msg"`
	Result     interface{}  `json:"result"`
}

func (c *MainController) StrategyList() (result []byte, code int32, err error) {
	// complete url
	strategyUrl := fmt.Sprintf("%s%s", model.GlobalConf.StrategyHost, model.GlobalConf.StrategyListUrl)

	// get jsonData struct
	jsonData := c.GetString("jsonData")
	var jsonDataStruct =  make(map[string]string)
	if err := json.Unmarshal([]byte(jsonData), &jsonDataStruct); err != nil {
		return nil, model.ParamsErrorCode, errors.New("jsonData type error")
	}

	// post strategy id list
	data := make(map[string]string)
	if model.GlobalConf.DefaultApiCode != "" {
		data["apiCode"] = model.GlobalConf.DefaultApiCode
	} else {
		data["apiCode"] = c.GetString("apiCode")
	}
	result, err = getData(strategyUrl, data)
	if err != nil {
		return result, model.ServerErrorCode, err
	}

	// result success judge
	strategyRsp := &StrategyListResponse{}
	err = json.Unmarshal(result, strategyRsp)
	if strategyRsp.Code != "00" {
		return nil, model.QueryErrorCode, errors.New("query strategy list fail")
	}
	response := model.Response{
		Code:    model.SuccessCode,
		Message: "success",
		Result:  strategyRsp.Data,
	}
	responseByte, _ := json.Marshal(response)

	return responseByte, model.SuccessCode, nil
}

func (c *MainController) StrategyStart() (result []byte, code int32, err error) {
	// complete url
	strategyUrl := fmt.Sprintf("%s%s", model.GlobalConf.StrategyHost, model.GlobalConf.StrategyStartUrl)

	// get jsonData struct
	jsonData := c.GetString("jsonData")
	var jsonDataStruct =  make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &jsonDataStruct); err != nil {
		return nil, model.ParamsErrorCode, errors.New("jsonData type error")
	}

	// get flow_id params
	var flow_id string
	var params map[string]interface{}
	var ok bool
	if _, ok := jsonDataStruct["flow_id"]; !ok {
		return nil, model.ParamsErrorCode, errors.New("no flow_id")
	}
	if flow_id, ok = jsonDataStruct["flow_id"].(string); !ok {
		return nil, model.ParamsErrorCode, errors.New("flow_id type error")
	}
	if _, ok := jsonDataStruct["params"]; !ok {
		return nil, model.ParamsErrorCode, errors.New("no parsms")
	}
	if params, ok = jsonDataStruct["params"].(map[string]interface{}); !ok {
		return nil, model.ParamsErrorCode, errors.New("params type error")
	}

	// gen flow_task_id and swift number
	flowTaskId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	var swiftNumber string
	if model.GlobalConf.DefaultApiCode != "" {
		swiftNumber = GenSwiftNumber(model.GlobalConf.DefaultApiCode)
	} else {
		swiftNumber = GenSwiftNumber(c.GetString("apiCode"))
	}

	// save cache
	if model.GlobalConf.NeedCache(flow_id) {
		var xmlData string
		if _, ok := jsonDataStruct["xml"]; !ok {
			return nil, model.ParamsErrorCode, errors.New("no xml data")
		}
		if xmlData, ok = jsonDataStruct["xml"].(string); !ok {
			return nil, model.ParamsErrorCode, errors.New("xml type error")
		}
		xmlData = strings.Replace(xmlData, "GBK", "UTF8", 1)
		paramsBytes, _ := json.Marshal(params)
		key := utils.GetMd5(string(paramsBytes))
		xmlReader := strings.NewReader(xmlData)
		jsonData, err := xj.Convert(xmlReader)
		if err != nil {
			return nil, model.ServerErrorCode, errors.New("xml parse error")
		}
		err = model.RedisClient.SetKeyExpire(key, jsonData.String(), 3600 * 24 * 7)
		if err != nil {
			return nil, model.ServerErrorCode, err
		}
	}

	// post strategy flow
	data := make(map[string]interface{})
	if model.GlobalConf.DefaultApiCode != "" {
		data["apiCode"] = model.GlobalConf.DefaultApiCode
	} else {
		data["apiCode"] = c.GetString("apiCode")
	}
	data["flow_id"] = flow_id
	data["flow_task_id"] = flowTaskId
	data["type"] = "1"
	data["params"] = params
	data["_platform"] = map[string]string{"appId": swiftNumber}

	result, err = postData(strategyUrl, data, true)
	if err != nil {
		return result, model.ServerErrorCode, err
	}

	// result success judge
	strategyRsp := &StrategyFlowResponse{}
	err = json.Unmarshal(result, strategyRsp)
	if strategyRsp.Code != "00" {
		return nil, model.QueryErrorCode, errors.New("query strategy flow fail")
	}
	response := model.Response{
		Code:    model.SuccessCode,
		Message: "success",
		Result:  strategyRsp.Result,
	}
	responseByte, _ := json.Marshal(response)

	return responseByte, model.SuccessCode, nil
}

