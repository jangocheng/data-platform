package composive

import (
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"platform/common/utils"
	"reflect"
	"strconv"
	"strings"
)

func GetMicroService(cnf interface{}, opts ...micro.Option) (micro.Service, error) {
	var flagOptions []cli.Flag
	for i := 0; i < reflect.ValueOf(cnf).Elem().NumField(); i++ {
		fieldName := utils.Camel2Snake(reflect.TypeOf(cnf).Elem().Field(i).Name)
		if fieldName == "registry" || fieldName == "registry_address" {
			continue
		}
		fieldType := reflect.TypeOf(cnf).Elem().Field(i).Type.String()
		usageTag := reflect.TypeOf(cnf).Elem().Field(i).Tag.Get("usage")
		defaultTag := reflect.TypeOf(cnf).Elem().Field(i).Tag.Get("default")
		if fieldType == "string" {
			flagOptions = append(flagOptions, cli.StringFlag{Name: fieldName, Usage: usageTag, Value: defaultTag})
		} else if fieldType == "bool" {
			flagOptions = append(flagOptions, cli.BoolFlag{Name: fieldName, Usage: usageTag})
		} else if strings.Contains(fieldType, "int") {
			flagOptions = append(flagOptions, cli.StringFlag{Name: fieldName, Usage: usageTag, Value: defaultTag})
		} else if strings.Contains(fieldType, "float") {
			flagOptions = append(flagOptions, cli.StringFlag{Name: fieldName, Usage: usageTag, Value: defaultTag})
		} else {
			panic(fmt.Sprintf("fail to create service receive wrong type config field: %s type: %s", fieldName, fieldType))
		}
	}
	flagOptions = append(flagOptions, cli.StringFlag{Name: "f", Usage: "provide config by file"})
	service := micro.NewService(append(opts, micro.Flags(flagOptions...))...)
	return service, nil
}

func InitService(service micro.Service, cnf interface{}, opts ...micro.Option) error {
	opts = append(opts, micro.Action(func(context *cli.Context) {
		for i := 0; i < reflect.ValueOf(cnf).Elem().NumField(); i++ {
			fieldName := reflect.TypeOf(cnf).Elem().Field(i).Name
			snakeFieldName := utils.Camel2Snake(fieldName)
			fieldType := reflect.TypeOf(cnf).Elem().Field(i).Type.String()
			switch fieldType {
			case "string":
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(context.String(snakeFieldName)))
			case "bool":
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(context.Bool(snakeFieldName)))
			case "int":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				intValue, err := strconv.Atoi(stringValue)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(intValue))
			case "uint":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				uintValue, err := strconv.ParseUint(stringValue, 10, 32)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(uint(uintValue)))
			case "uint8":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				uint8Value, err := strconv.ParseUint(stringValue, 10, 8)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(uint8(uint8Value)))
			case "uint16":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				uint16Value, err := strconv.ParseUint(stringValue, 10, 16)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(uint16(uint16Value)))
			case "uint32":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				uint32Value, err := strconv.ParseUint(stringValue, 10, 32)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(uint32(uint32Value)))
			case "uint64":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				uint64Value, err := strconv.ParseUint(stringValue, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(uint64Value))
			case "int8":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				int8Value, err := strconv.ParseInt(stringValue, 10, 8)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(int8(int8Value)))
			case "int16":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				int16Value, err := strconv.ParseInt(stringValue, 10, 16)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(int16(int16Value)))
			case "int32":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				int32Value, err := strconv.ParseInt(stringValue, 10, 32)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(int32(int32Value)))
			case "int64":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				int64Value, err := strconv.ParseInt(stringValue, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(int64Value))
			case "float32":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				float32Value, err := strconv.ParseFloat(stringValue,  32)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(float32(float32Value)))
			case "float64":
				stringValue := context.String(snakeFieldName)
				if stringValue == "" {
					stringValue = "0"
				}
				float64Value, err := strconv.ParseFloat(stringValue,  64)
				if err != nil {
					panic(fmt.Sprintf("fail to init service %s", err))
				}
				reflect.ValueOf(cnf).Elem().FieldByName(fieldName).Set(reflect.ValueOf(float64Value))
			default:
				panic(fmt.Sprintf("receive wrong field type config field of name: %s type:%s", fieldName, fieldType))
			}
		}
		if context.String("f") != "" {
			err := utils.LoadJsonConf(context.String("f"), cnf)
			if err != nil {
				panic(fmt.Sprintf("fail to init service %s", err))
			} else {
				return
			}
		}
	}))
	service.Init(opts...)
	return nil
}
