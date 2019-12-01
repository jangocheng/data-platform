package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"math/rand"
	"net/url"
	HttpRequest "platform/common/http"
	"platform/data-proxy/model"
	"strings"
	"time"
)


type MainController struct {
	beego.Controller
}

func (c *MainController) RouterDefine(data map[string]string) func() ([]byte, int32, error) {
	prefixRouterMap := make(map[string]func() ([]byte, int32, error))

	prefixRouterMap = map[string]func() ([]byte, int32, error){
		"varset.biz": c.Bairong,
		"varset.hx": c.Bairong,
		"varset.platform": c.Platform,
		"var.platform": c.Platform,
		"StrategyFlowStart": c.StrategyStart,
		"StrategyFlowList": c.StrategyList,
		"PostLoanList": c.PostLoanList,
		"PostLoanStart": c.PostLoanStart,
	}

	for k, v := range prefixRouterMap {
		if strings.HasPrefix(data["api"], k) {
			return v
		}
	}

	for k, v := range prefixRouterMap {
		if strings.HasPrefix(data["varset"], k) {
			return v
		}
	}

	for k, v := range prefixRouterMap {
		if strings.HasPrefix(data["var"], k) {
			return v
		}
	}

	return nil
}

func (c *MainController) RouterComplete(status int32, msg string, startTime time.Time,
							runtime model.Runtime, result ...[]byte) {

	runtime.FormData = getFormMap(c.Ctx.Request.Form)
	runtime.RequestTime = startTime.Format("2006-01-02 15:04:05")
	runtime.Message = msg
	runtime.CostTime = int32(time.Now().Sub(startTime).Milliseconds())
	runtime.Status = status

	_ = model.Logger.Log(runtime)

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	if result != nil && len(result) != 0 && status == model.SuccessCode {
		_, err := c.Ctx.ResponseWriter.Write(result[0])
		if err != nil {
			status = model.ServerErrorCode
			msg = fmt.Sprintf("fail to set data to response error info: %s", err.Error())
		} else {
			return
		}
	}

	c.Data["json"] = model.Response{
		Code:  		status,
		Message: 	msg,
	}
	c.ServeJSON()
}


// apiCode
// jsonData

func (c *MainController) Router() {
	runtime := model.Runtime{}
	startTime := time.Now()

	// verify url correct
	if _, ok := model.PermissionUrlMap[c.Ctx.Request.URL.Path]; !ok {
		errMsg := "url invalid"
		c.RouterComplete(model.InvalidApiCode, errMsg, startTime, runtime)
		return
	}

	// verify apiCode
	apiCode := c.GetString("apiCode")
	if apiCode == "" {
		c.RouterComplete(model.ParamsErrorCode, "apiCode null", startTime, runtime)
		return
	}

	// verify prod code
	codeMap := make(map[string]string)
	jsonData := make(map[string]interface{})
	_ = json.Unmarshal([]byte(c.GetString("jsonData")), &jsonData)
	if jsonData["api"] != nil {
		codeMap["api"] = jsonData["api"].(string)
	} else if jsonData["varset"] != nil {
		codeMap["varset"] = jsonData["varset"].(string)
	} else if jsonData["var"] != nil {
		codeMap["var"] = jsonData["var"].(string)
	}else {
		errMsg := "no prod code"
		c.RouterComplete(model.ParamsErrorCode, errMsg, startTime, runtime)
		return
	}

	// get handler func
	exeFunc := c.RouterDefine(codeMap)
	if exeFunc == nil {
		errMsg := "can not find prod code exec func"
		c.RouterComplete(model.ParamsErrorCode, errMsg, startTime, runtime)
		return
	}

	// get handler  result
	result, code, err := exeFunc()
	if err != nil {
		c.RouterComplete(code, err.Error(), startTime, runtime)
		return
	}
	if result == nil || len(result) == 0 {
		c.RouterComplete(model.NoResultCode, "handler return nil", startTime, runtime)
		return
	}
	c.RouterComplete(code, "success", startTime, runtime, result)
}

func getFormMap(form url.Values) map[string]string {
	formMap := make(map[string]string)
	for k, v := range form {
		if len(v) > 0 {
			formMap[k] = v[0]
		} else {
			formMap[k] = ""
		}
	}
	return formMap
}

func postData(url string, data interface{}, json bool) (result []byte, err error) {

	reqCon := HttpRequest.RequestConfig{
		Url: url,
		Data: data,
		Timeout: 60,
		DisTlsVerify:  true,
	}
	var response *HttpRequest.Response
	if json == true {
		response, err = HttpRequest.NewRequest(reqCon).PostJson()
	} else {
		response, err = HttpRequest.NewRequest(reqCon).Post()
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to get response")
	}

	body, err := response.Body()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get response body")
	}
	return body, nil
}

func getData(url string, params map[string]string) (result []byte, err error) {

	reqCon := HttpRequest.RequestConfig{
		Url: url,
		Params: params,
		Timeout: 60,
		DisTlsVerify:  true,
	}
	response, err := HttpRequest.NewRequest(reqCon).Get()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get response")
	}

	body, err := response.Body()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get response body")
	}
	return body, nil
}


func GenSwiftNumber(apiCode string) string {
	if apiCode == "" {
		apiCode = "0000000"
	}
	swiftNumber := fmt.Sprintf("%s_%s_%d", apiCode, time.Now().Format("20060102150405060"), rand.Intn(999999))
	return swiftNumber
}
