package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func GenSwiftNumber(paramsStr string) string {
	var apiCode string
	var swiftNumber string
	params := make(map[string]interface{})
	_ = json.Unmarshal([]byte(paramsStr), &params)

	if params["apiCode"] != nil {
		apiCode = fmt.Sprintf("%v", params["apiCode"])
	} else {
		apiCode = "0000000"
	}
	swiftNumber = fmt.Sprintf("%s_%s_%d", apiCode, time.Now().Format("20060102150405060"), rand.Intn(999999))
	return swiftNumber
}


func GetParentSwiftNumber(params map[string]string) string {
	var hiddenParams map[string]string
	hiddenParams = make(map[string]string)
	if platform, ok := params["_platform"]; ok {
		if err := json.Unmarshal([]byte(platform), &hiddenParams); err == nil {
			return hiddenParams["parent_swift_number"]
		} else {
			return ""
		}
	} else {
		return ""
	}
}



