package otto

import (
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
	"platform/common/utils"
	cvm "platform/common/vm"
	"time"
)

func setFunction(name string, function interface{}, vm cvm.ClientVM) error {
	err := vm.SetFunc(name, function)
	if err != nil {
		return err
	}
	return nil
}

func SetFunctionTools(cvm cvm.ClientVM) error {
	err := setFunction("md5", func(call otto.FunctionCall) otto.Value {
		value, err := otto.ToValue(utils.GetMd5(call.Argument(0).String()))
		if err != nil {
			value, _ = otto.ToValue("")

			return value
		}
		return value
	},
	cvm,
	)
	if err != nil {
		return errors.Wrap(err, "fail to set func tools to js vm otto")
	}
	err = setFunction("getPreTime", func(call otto.FunctionCall) otto.Value {
		value, err :=  otto.ToValue(utils.GetPreTime(call.Argument(0).String()))
		if err != nil {
			value, _ = otto.ToValue("")

			return value
		}
		return value
	},
	cvm,
	)
	if err != nil {
		return errors.Wrap(err, "fail to set func tools to js vm otto")
	}
	err = setFunction("timeFormat", func(call otto.FunctionCall) otto.Value {
		value, err := otto.ToValue(utils.TimeFormat(call.Argument(0).String()))
		if err != nil {
			value, _ = otto.ToValue("")
			return value
		}
		return value
	},
		cvm,
	)
	if err != nil {
		return errors.Wrap(err, "fail to set func tools to js vm otto")
	}
	err = setFunction("timestamp", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			value, err := otto.ToValue(utils.Timestamp())
			if err != nil {
				value, _ = otto.ToValue(time.Now().Unix())
				return value
			}
			return value
		} else {
			value, err := otto.ToValue(utils.Timestamp(call.Argument(0).String()))
			if err != nil {
				value, _ = otto.ToValue(time.Now().Unix())
				return value
			}
			return value
		}
	},
		cvm,
	)
	if err != nil {
		return errors.Wrap(err, "fail to set func tools to js vm otto")
	}
	return nil
}
