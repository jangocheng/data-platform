package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"platform/data-proxy/model"
)

type PostLoanListResponse struct {
	Code       string  		`json:"code"`
	Msg        string   	`json:"msg"`
	Data       interface{}  `json:"data"`
}

type PostLoanStartResponse struct {
	Code       string  		`json:"code"`
	Msg        string   	`json:"msg"`
}

func (c *MainController) PostLoanList() (result []byte, code int32, err error) {
	// complete url
	postLoanUrl := fmt.Sprintf("%s%s", model.GlobalConf.PostLoanHost, model.GlobalConf.PostLoanListUrl)

	// get jsonData struct
	jsonData := c.GetString("jsonData")
	var jsonDataStruct =  make(map[string]string)
	if err := json.Unmarshal([]byte(jsonData), &jsonDataStruct); err != nil {
		return nil, model.ParamsErrorCode, errors.New("jsonData type error")
	}

	// post postloan id list
	data := make(map[string]string)
	if model.GlobalConf.DefaultApiCode != "" {
		data["apiCode"] = model.GlobalConf.DefaultApiCode
	} else {
		data["apiCode"] = c.GetString("apiCode")
	}
	result, err = getData(postLoanUrl, data)
	if err != nil {
		return result, model.ServerErrorCode, err
	}

	// result success judge
	postLoanRsp := &PostLoanListResponse{}
	err = json.Unmarshal(result, postLoanRsp)
	if postLoanRsp.Code != "00" {
		return nil, model.QueryErrorCode, errors.New("query postLoad list fail")
	}
	response := model.Response{
		Code:    model.SuccessCode,
		Message: "success",
		Result:  postLoanRsp.Data,
	}
	responseByte, _ := json.Marshal(response)
	return responseByte, model.SuccessCode, nil
}

func (c *MainController) PostLoanStart() (result []byte, code int32, err error) {
	// complete url
	postLoanUrl := fmt.Sprintf("%s%s", model.GlobalConf.PostLoanHost, model.GlobalConf.PostLoanStartUrl)

	// get jsonData struct
	jsonData := c.GetString("jsonData")
	var jsonDataStruct =  make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &jsonDataStruct); err != nil {
		return nil, model.ParamsErrorCode, errors.New("jsonData type error")
	}

	// get task_id params
	var task_id string
	var entityList []interface{}
	var ok bool
	if _, ok := jsonDataStruct["taskId"]; !ok {
		return nil, model.ParamsErrorCode, errors.New("no task_id")
	}
	if task_id, ok = jsonDataStruct["taskId"].(string); !ok {
		return nil, model.ParamsErrorCode, errors.New("task_id type error")
	}
	if _, ok := jsonDataStruct["entityList"]; !ok {
		return nil, model.ParamsErrorCode, errors.New("no parsms")
	}
	if entityList, ok = jsonDataStruct["entityList"].([]interface{}); !ok {
		return nil, model.ParamsErrorCode, errors.New("entityList type error")
	}

	// gen swift number
	var swiftNumber string
	if model.GlobalConf.DefaultApiCode != "" {
		swiftNumber = GenSwiftNumber(model.GlobalConf.DefaultApiCode)
	} else {
		swiftNumber = GenSwiftNumber(c.GetString("apiCode"))
	}

	// post postLoan start
	data := make(map[string]interface{})

	if model.GlobalConf.DefaultApiCode != "" {
		data["apiCode"] = model.GlobalConf.DefaultApiCode
	} else {
		data["apiCode"] = c.GetString("apiCode")
	}
	data["taskId"] = task_id
	data["entityList"] = entityList
	data["_platform"] = map[string]string{"appId": swiftNumber}

	result, err = postData(postLoanUrl, data, true)
	if err != nil {
		return result, model.ServerErrorCode, err
	}

	// result success judge
	postLoanRsp := &PostLoanStartResponse{}
	err = json.Unmarshal(result, postLoanRsp)
	if postLoanRsp.Code != "00" {
		return nil, model.QueryErrorCode, errors.New("query postLoan result fail")
	}
	response := model.Response{
		Code:    model.SuccessCode,
		Message: "success",
	}
	responseByte, _ := json.Marshal(response)

	return responseByte, model.SuccessCode, nil


}
