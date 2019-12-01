package utils

import "strings"

func SplitValString(str string) []string {
	var resultList []string
	resultList = strings.Split(str, ",")
	for _, i := range resultList {
		if strings.Trim(i, " ") == "" {
			continue
		} else {
			resultList = append(resultList, strings.Trim(i, " "))
		}
	}
	return resultList
}
