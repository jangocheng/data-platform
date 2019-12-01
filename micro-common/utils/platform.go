package utils

import "encoding/json"

func GetPlatformParam(parentSwiftNumber string) string {
	params := make(map[string]string)
	params["parent_swift_number"] = parentSwiftNumber
	jsonPlatform, _ := json.Marshal(params)
	return string(jsonPlatform)




}
