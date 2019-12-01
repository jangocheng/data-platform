package mysql

import (
	"reflect"
	"strings"
	"time"
)

func checkMustFieldExist(ptr interface{}) (bool, string) {
	reflectType := reflect.TypeOf(ptr).Elem()
	reflectValue := reflect.ValueOf(ptr).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		xormTag := reflectType.Field(i).Tag.Get("xorm")
		if strings.Contains(xormTag, "notnull") &&
			!strings.Contains(xormTag, "default") &&
			!strings.Contains(xormTag, "pk") {
			switch reflectValue.Field(i).Interface().(type) {
			case time.Time:
				value, _ := reflectValue.Field(i).Interface().(time.Time)
				nulltime := time.Time{}
				if value == nulltime {
					return false, reflectType.Field(i).Name
				}
			case int64:
				value, _ := reflectValue.Field(i).Interface().(int64)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case int:
				value, _ := reflectValue.Field(i).Interface().(int)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case int32:
				value, _ := reflectValue.Field(i).Interface().(int32)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case int8:
				value, _ := reflectValue.Field(i).Interface().(int8)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case uint:
				value, _ := reflectValue.Field(i).Interface().(uint)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case uint8:
				value, _ := reflectValue.Field(i).Interface().(uint8)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case uint32:
				value, _ := reflectValue.Field(i).Interface().(uint32)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case uint64:
				value, _ := reflectValue.Field(i).Interface().(uint64)
				if value == 0 {
					return false, reflectType.Field(i).Name
				}
			case string:
				value, _ := reflectValue.Field(i).Interface().(string)
				if value == "" {
					return false, reflectType.Field(i).Name
				}
			case []byte:
				value, _ := reflectValue.Field(i).Interface().([]byte)
				if value == nil {
					return false, reflectType.Field(i).Name
				}
			}
		}

	}
	return true, ""
}

func mustFieldNotSetVerify(ptr interface{}) (bool, string) {
	sliceValue := reflect.Indirect(reflect.ValueOf(ptr))
	if sliceValue.Kind() == reflect.Slice {
		size := sliceValue.Len()
		for i := 0; i < size; i++ {
			ptr = sliceValue.Index(i).Interface()
			result, msg := checkMustFieldExist(ptr)
			if !result {
				return result, msg
			}
		}
		return true, ""
	} else {
		return checkMustFieldExist(ptr)
	}
}
