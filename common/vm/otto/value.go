package otto

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
	"platform/common/vm"
)


type ValueOtto struct {
	value    otto.Value
}


func (v *ValueOtto) GetString() string {
	return v.value.String()
}

func (v *ValueOtto) GetFloat() (float64, error) {
	floatValue, err := v.value.ToFloat()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get float from value")
	}
	return floatValue, nil
}

func (v *ValueOtto) GetObjectList() ([]vm.Value, error) {
	var resultList []vm.Value
	keys := v.value.Object().Keys()
	for _, key := range keys {
		keyValue, err := v.value.Object().Get(key)
		if err != nil {
			return nil, errors.Wrap(err, "fail to get object list from value")
		}
		resultList = append(resultList, &ValueOtto{value:keyValue})
	}
	return resultList, nil
}

func (v *ValueOtto) GetObjectKey(key string) (vm.Value, error) {
	keyValue, err := v.value.Object().Get(key)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to get object map key %s from value", key))
	}
	return &ValueOtto{value:keyValue}, nil
}

func (v *ValueOtto) GetInt() (int64, error) {
	intValue, err := v.value.ToInteger()
	if err != nil {
		return 0, errors.Wrap(err, "fail to get int from value")
	}
	return intValue, nil
}

func (v *ValueOtto) GetBool() (bool, error) {
	boolValue, err := v.value.ToBoolean()
	if err != nil {
		return false, errors.Wrap(err, "fail to get bool from value")
	}
	return boolValue, nil
}


func (v *ValueOtto) IsNil() bool {
	if v.value.IsNaN() || v.value.IsNull() {
		return true
	}
	return false
}

