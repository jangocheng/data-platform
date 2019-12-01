package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	eMap "platform/common/expire-map"
	"platform/data-proxy/model"
)

var (
	expireMap = eMap.NewExpiredMap()
	tokenId   = "TOKEN"
)


func (c *MainController) Bairong() (result []byte, code int32, err error) {

	token, err := c.GetBairongToken()
	if err != nil {
		return nil, model.SuccessCode, err
	}

	reqForm := getFormMap(c.Ctx.Request.Form)
	reqForm["tokenid"] = token

	jsonData := make(map[string]interface{})
	_ = json.Unmarshal([]byte(c.GetString("jsonData")), &jsonData)

	if jsonData["api"] != nil {
		url := fmt.Sprintf("%s%s", model.GlobalConf.BairongBizHost,
			model.GlobalConf.BairongBizApiPath)
		result, err = postData(url, reqForm, false)

	} else if jsonData["derive"] != nil {
		url := fmt.Sprintf("%s%s", model.GlobalConf.BairongBizHost,
			model.GlobalConf.BairongBizDerivePath)
		result, err = postData(url, reqForm, false)
	}
	return result, model.SuccessCode, err
}


func (c *MainController) GetBairongToken() (string, error) {
	found, tokenInterface := expireMap.Get(tokenId)
	if found {
		return tokenInterface.(string), nil
	}

	apiData := make(map[string]string)
	if model.GlobalConf.DefaultApiCode != "" {
		apiData["apiCode"] = model.GlobalConf.DefaultApiCode
	} else {
		apiData["apiCode"] = c.GetString("apiCode")
	}
	apiData["userName"] = model.GlobalConf.BairongUserName
	apiData["password"] = model.GlobalConf.BairongPassword

	url := fmt.Sprintf("%s%s", model.GlobalConf.BairongHost,
		model.GlobalConf.BairongLoginPath)

	body, err := postData(url, apiData, false)

	loginRsp := make(map[string]string)
	err = json.Unmarshal(body, &loginRsp)
	if err != nil {
		return "", errors.Wrap(err, "transfer bairong login token body fail")
	}
	if loginRsp["code"] != "00" {
		return "", errors.Wrap(err, fmt.Sprintf("bairong login response wrong status: %s", loginRsp["code"]))
	}
	tokenStr := loginRsp["tokenid"]
	expireMap.Set(tokenId, tokenStr, 60 * 60 * 1.9)
	return tokenStr, nil

}


