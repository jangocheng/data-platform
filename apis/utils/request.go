package utils

import api "github.com/micro/go-micro/api/proto"

func GetMethodParams(req *api.Request) (map[string]string, map[string]string, map[string]string) {
	var getParams map[string]string
	var postParams map[string]string
	var allParams map[string]string
	getParams = make(map[string]string, 0)
	postParams = make(map[string]string, 0)
	allParams = make(map[string]string, 0)
	for key, values := range req.Get {
		getParams[key] = values.Values[0]
		allParams[key] = values.Values[0]
	}
	for key, values := range req.Post {
		postParams[key] = values.Values[0]
		allParams[key] = values.Values[0]
	}
	return getParams, postParams, allParams
}
