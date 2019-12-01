package data_source

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func SplitParamsStr(paramsStr string) (params map[string]string, must []string) {
	params = make(map[string]string)
	paramsNeed := strings.Split(paramsStr, ",")
	for _, paramStr := range paramsNeed {
		if paramStr == "" {
			continue
		}
		var paramNameDefault []string
		if strings.Contains(paramStr, ":") {
			paramNameDefault = strings.Split(paramStr, ":")
		} else if strings.Contains(paramStr, "=") {
			paramNameDefault = strings.Split(paramStr, "=")
		} else {
			paramNameDefault = []string{paramStr}
			must = append(must, paramStr)
		}
		if len(paramNameDefault) > 1 {
			params[paramNameDefault[0]] = paramNameDefault[1]
		} else {
			params[paramNameDefault[0]] = ""
		}
	}
	return
}

func MustParamsRecVerify(paramsStr string, reqParams map[string]interface{}) error {
	paramNotRecStr := ""
	_, must := SplitParamsStr(paramsStr)
	for _, paramMust := range must {
		if _, ok := reqParams[paramMust]; !ok {
			minParams := paramMust + "_min"
			maxParams := paramMust + "_max"
			if _, ok := reqParams[minParams]; ok {
				continue
			}
			if _, ok := reqParams[maxParams]; ok {
				continue
			}
			paramNotRecStr += fmt.Sprintf("%s ", paramMust)
		}
	}
	if paramNotRecStr == "" {
		return nil
	}
	return errors.New(fmt.Sprintf("param %smust be sent", paramNotRecStr))
}


func GetMysqlCondition(paramsNeedStr string, reqParams map[string]interface{}) (string, []interface{}) {
	whereStr := ""
	var args []interface{}
	dsParams, _ := SplitParamsStr(paramsNeedStr)

	for name, defaultV := range dsParams {
		if value, ok := reqParams[name]; ok {
			if value == "" || value == nil {
				continue
			}
			whereStr += name + "=" + "?" + " and "
			args = append(args, value)
		} else {
			if defaultV == "" {
				if rangeParams, ok := reqParams[name + "_min"]; ok {
					whereStr += name + ">=" + "?" + " and "
					args = append(args, rangeParams)
				}
				if rangeParams, ok := reqParams[name + "_max"]; ok {
					whereStr += name + "<=" + "?" + " and "
					args = append(args, rangeParams)
				}
			} else {
				whereStr += name + "=" + "?" + " and "
				args = append(args, defaultV)
			}
		}
	}
	whereStr = strings.TrimRight(whereStr, " and ")

	return whereStr, args
}

func GetOracleCondition(paramsNeedStr string, reqParams map[string]interface{}) (string, []interface{}) {
	whereStr := ""
	var args []interface{}
	dsParams, _ := SplitParamsStr(paramsNeedStr)

	for name, defaultV := range dsParams {
		if value, ok := reqParams[name]; ok {
			if value == "" || value == nil {
				continue
			}
			whereStr += "\"" + name + "\"" + "=" + "?" + " and "
			args = append(args, value)
		} else {
			if defaultV == "" {
				if rangeParams, ok := reqParams[name + "_min"]; ok {
					whereStr += "\"" + name + "\"" + ">=" + "?" + " and "
					args = append(args, rangeParams)
				}
				if rangeParams, ok := reqParams[name + "_max"]; ok {
					whereStr += "\"" + name + "\"" + "<=" + "?" + " and "
					args = append(args, rangeParams)
				}
			} else {
				whereStr += "\"" + name + "\"" + "=" + "?" + " and "
				args = append(args, defaultV)
			}
		}
	}
	whereStr = strings.TrimRight(whereStr, " and ")

	return whereStr, args
}

func GetSqlStr(sqlRow string, paramsNeedStr string, reqParams map[string]interface{}) string {
	var sql = sqlRow
	dsParams, _ := SplitParamsStr(paramsNeedStr)
	for name, defaultV := range dsParams {
		if _, ok := reqParams[name]; ok {
			sql = strings.Replace(sql, fmt.Sprintf("{%s}", name), fmt.Sprintf("%v", reqParams[name]), -1)
		} else {
			sql = strings.Replace(sql, fmt.Sprintf("{%s}", name), fmt.Sprintf("%v", defaultV), -1)
		}
	}
	return sql
}

func GetMongoCondition(paramsNeedStr string, reqParams map[string]interface{}) map[string]interface{} {
	condition := make(map[string]interface{}, 0)
	dsParams, _ := SplitParamsStr(paramsNeedStr)

	for name, defaultV := range dsParams {
		if value, ok := reqParams[name]; ok {
			if value == "" {
				continue
			}
			condition[name] = value
		} else if rangeV, ok := reqParams[name + "_min"]; ok {
			rangeValue := make(map[string]interface{})
			rangeValue["$gte"] = rangeV
			condition[name] = rangeValue
		} else if rangeV, ok := reqParams[name + "_max"]; ok {
			rangeValue := make(map[string]interface{})
			rangeValue["$lte"] = rangeV
			condition[name] = rangeValue
		} else {
			if defaultV == "" {
				continue
			}
			condition[name] = defaultV
		}
	}
	return condition
}