package data_api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func (d *DataApiWeb) JsWrapper(wrappers []Wrapper, params interface{}) error {
	for _, wrapper := range wrappers {
		jsCode := wrapper.JsCode
		_, err := d.ottoClient.Run(jsCode, params)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("fail to execute js %s wrapper: %s", wrapper.WrapperCode, err.Error()))
		}
	}
	return nil
}

func (d *DataApiWeb) JsCustomWrapper(wrappers []Wrapper, params *AfterWrapperMsg) error {
	for _, wrapper := range wrappers {
		if params.Fail == true {
			return nil
		}
		jsCode := wrapper.JsCode
		_, err := d.ottoClient.Run(jsCode, params)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("fail to execute js %s wrapper: %s", wrapper.WrapperCode, err.Error()))
		}
	}
	if len(wrappers) != 0 {
		_, err := d.ottoClient.Run(`input.Result = JSON.stringify(input)`, params)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("fail to execute js wrapper: %s", err.Error()))
		}
		err = json.Unmarshal([]byte(params.Result.(string)), params)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("fail to execute js wrapper: %s", err.Error()))
		}
	}
	return nil
}


func (d *DataApiWeb) LogWrapper(wrappers []Wrapper, logModel *LogWrapperMsg) int32 {
	var executeSuccessTime int32
	for _, wrapper := range wrappers {
		if logModel.Fail == true {
			return 0
		}
		jsCode := wrapper.JsCode
		_, err := d.ottoClient.Run(jsCode, logModel)
		if err == nil {
			executeSuccessTime += 1
		} else {
			return executeSuccessTime
		}
	}

	if len(wrappers) != 0 {
		_, err := d.ottoClient.Run(`input.Result = JSON.stringify(input)`, logModel)
		if err != nil {
			return executeSuccessTime
		}
		err = json.Unmarshal([]byte(logModel.Result.(string)), logModel)
		if err != nil {
			return executeSuccessTime
		}
	}
	return executeSuccessTime
}


