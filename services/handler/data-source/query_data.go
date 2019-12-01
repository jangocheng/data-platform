package data_source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/util/file"
	"github.com/pkg/errors"
	"path"
	"platform/common/ftp"
	"platform/common/http"
	"platform/common/sql"
	"platform/common/utils"
	dataSource "platform/proto/service/data-source"
	"reflect"
	"strconv"
	"strings"
	"time"
)
var (
	HttpWrongStatusError = errors.New("wrong http response status")
)

func (d *DataSourceService) QueryData(ctx context.Context, req *dataSource.QueryDataRequest, rsp *dataSource.QueryDataResponse) error {
	startTime := time.Now()
	if req.Params == "" {
		req.Params = "{}"
	}
	dsType := req.DsType
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.Params), &params)
	if err != nil {
		rsp.Status = ParamsErrorStatus
		rsp.Message = "request Param params format wrong"
		goto Over
	}
	if req.DsCode == "" || req.DsType == "" {
		rsp.Status = ParamsErrorStatus
		rsp.Message = "request params ds type and ds code can not be none"
		goto Over
	}
	switch dsType {
	case "file":
		paramsStr := make(map[string]string)
		err := json.Unmarshal([]byte(req.Params), &params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = "request Param param format wrong"
			break
		}
		var item FileDataSource
		var got bool
		if item, got = fileDataSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		if _, ok := dataOwnerMapInstance[item.OrgCode]; !ok {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "org code not found"
			break
		}
		result, err := d.queryFileData(item, paramsStr["file"])
		if err != nil {
			rsp.Status = ErrorStatus
			rsp.Message = fmt.Sprintf("fail to get file data %s", err.Error())
			break
		}
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		rsp.Result = result
		break
	case "http":
		paramsStr := make(map[string]string)
		err := json.Unmarshal([]byte(req.Params), &params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = "request Param param format wrong"
			break
		}
		var item HttpDataSource
		var got bool
		if item, got = httpDataSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		if _, ok := dataOwnerMapInstance[item.OrgCode]; !ok {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "org code not found"
			break
		}
		err = MustParamsRecVerify(item.Params, params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = err.Error()
			break
		}
		result, err := d.queryHttpData(item, paramsStr, req.RawData)
		if err != nil {
			if errors.Cause(err) == HttpWrongStatusError {
				rsp.Status = HttpWrongStatus
				rsp.Message = err.Error()
				rsp.Result = result
				break
			}
			rsp.Status = ErrorStatus
			rsp.Message = fmt.Sprintf("fail to get http data %s", err.Error())
			break
		}
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		rsp.Result = result
		break
	case "db":
		var item DatabaseSource
		var got bool
		if item, got = databaseSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		if _, ok := dataOwnerMapInstance[item.OrgCode]; !ok {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "org code not found"
			break
		}
		err = MustParamsRecVerify(item.Params, params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = err.Error()
			break
		}
		result, err := d.queryDBData(item, params)
		if err != nil {
			rsp.Status = ErrorStatus
			rsp.Message = fmt.Sprintf("fail to get database data %s", err.Error())
			break
		}
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		rsp.Result = result
		break
	case "kv":
		paramsStr := make(map[string]string)
		err := json.Unmarshal([]byte(req.Params), &params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = "request Param param format wrong"
			break
		}
		var item KVDatabaseSource
		var got bool
		if item, got = kvDatabaseSourceMapInstance[req.DsCode]; !got {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "ds code not found"
			break
		}
		if item.Invalid == 1 {
			rsp.Status = InvalidStatus
			rsp.Message = "data source is invalid"
			break
		}
		if _, ok := dataOwnerMapInstance[item.OrgCode]; !ok {
			rsp.Status = NotFundErrorStatus
			rsp.Message = "org code not found"
			break
		}
		err = MustParamsRecVerify("key", params)
		if err != nil {
			rsp.Status = ParamsErrorStatus
			rsp.Message = err.Error()
			break
		}
		result, err := d.queryKVData(item, paramsStr["key"])
		if err != nil {
			rsp.Status = ErrorStatus
			rsp.Message = fmt.Sprintf("fail to get http data %s", err.Error())
			break
		}
		rsp.Status = SuccessStatus
		rsp.Message = "success"
		rsp.Result = result
		break
	default:
		rsp.Status = ParamsErrorStatus
		rsp.Message = "ds type not found"
	}
	Over:
		costTime := int32(time.Now().Sub(startTime).Milliseconds())
		rsp.CostTime = costTime
		_ = d.logger.Log(logFormat{
			Status:            rsp.Status,
			Path:              reflect.TypeOf(req).Elem().String(),
			CostTime:          costTime,
			RequestTime:       startTime.Format("2006-01-02 15:04:05"),
			ParentSwiftNumber: req.XParentSwiftNumber,
			FormData:          req,
			Result:            rsp.Result,
			Message:           rsp.Message,
		})
	return nil
}

func (d *DataSourceService) queryFileData(item FileDataSource, fileNames string) (string, error) {
	var filesContent []string
	dataOwner := dataOwnerMapInstance[item.OrgCode]
	if dataOwner.Addr == "" {
		if fileNames == "" {
			fileList := utils.ReadDir(item.FileFolder)
			return  "[" + strings.Join(fileList, ", ") + "]", nil
		} else {
			for _, fileName := range strings.Split(fileNames, ",") {
				fileName := strings.Trim(fileName, " ")
				if fileName == "" {
					continue
				}
				fileFullPath := path.Join(item.FileFolder, fileName)
				exist, _ := file.Exists(fileFullPath)
				if !exist {
					return "", errors.New("file name not found")
				}
				content, err := utils.GetFileContent(fileFullPath)
				if err != nil {
					return "", err
				}
				filesContent = append(filesContent, content)
			}
		}
	} else {
		hostPortList := strings.Split(dataOwner.Addr, ":")
		if len(hostPortList) != 2 {
			return "", errors.New("the data owner addr is wrong")
		}
		host := hostPortList[0]
		port, err := strconv.Atoi(hostPortList[1])
		if err != nil {
			return "", errors.New("the data owner addr is wrong")
		}
		if fileNames == "" {
			fileList, err := ftp.ReadDir(dataOwner.User, dataOwner.PassWord, host, port, item.FileFolder)
			if err != nil {
				return "", err
			}
			return "[" + strings.Join(fileList, ", ") + "]", nil
		} else {

			for _, fileName := range strings.Split(fileNames, ",") {
				fileName := strings.Trim(fileName, " ")
				if fileName == "" {
					continue
				}
				fileFullPath := path.Join(item.FileFolder, fileName)
				content, err := ftp.ReadFile(dataOwner.User, dataOwner.PassWord, host, port, fileFullPath)
				if err != nil {
					return "", err
				}
				filesContent = append(filesContent, content)
			}
		}
	}
	filesContentStr, _ := json.Marshal(filesContent)

	return string(filesContentStr), nil
}

func (d *DataSourceService) queryHttpData(item HttpDataSource, params map[string]string, rawData string) (string, error) {
	headerMap := make(map[string]string)
	if params["_headers"] != "" {
		err := json.Unmarshal([]byte(params["_headers"]), &headerMap)
		if err != nil {
			return "", err
		}
	}

	cookieMap := make(map[string]string)
	if params["_cookies"] != "" {
		err := json.Unmarshal([]byte(params["_cookies"]), &cookieMap)
		if err != nil {
			return "", err
		}
	}
	switch strings.ToLower(item.Method) {
	case "get":
		domin := strings.TrimRight(dataOwnerMapInstance[item.OrgCode].Addr, "/")
		reqCon := HttpRequest.RequestConfig{
			Url:  domin + item.Path,
			Params: params,
			Timeout: 16,
			Cookies: cookieMap,
			Headers: headerMap,
			DisTlsVerify:  true,
		}
		request := HttpRequest.NewRequest(reqCon)
		response, err := request.Get()
		if err != nil {
			return "", err
		}
		rspBody, err := response.Body()
		if err != nil {
			return "", err
		}
		if strconv.Itoa(response.StatusCode()) != item.SuccessCode {
			return string(rspBody), errors.WithMessage(HttpWrongStatusError,
				fmt.Sprintf("wrong status: %d", response.StatusCode()))
		}
		return string(rspBody), nil
	case "post":
		domin := strings.TrimRight(dataOwnerMapInstance[item.OrgCode].Addr, "/")
		var reqCon HttpRequest.RequestConfig
		if rawData != "" {
			reqCon = HttpRequest.RequestConfig{
				Url:  domin + item.Path,
				Params: params,
				Data: rawData,
				Timeout: 16,
				Cookies: cookieMap,
				Headers: headerMap,
				DisTlsVerify:  true,
			}
		} else {
			reqCon = HttpRequest.RequestConfig{
				Url:  domin + item.Path,
				Data: params,
				Timeout: 16,
				Cookies: cookieMap,
				Headers: headerMap,
				DisTlsVerify:  true,
			}
		}
		request := HttpRequest.NewRequest(reqCon)
		var response *HttpRequest.Response
		var err error
		if item.ContentType == "json" {
			response, err = request.PostJson()
		} else if item.ContentType == "form" || item.ContentType == "" {
			response, err = request.Post()
		} else {
			err = errors.New("content type is not available")
		}
		if err != nil {
			return "", err
		}
		rspBody, err := response.Body()
		if err != nil {
			return "", err
		}
		if strconv.Itoa(response.StatusCode()) != item.SuccessCode {
			return string(rspBody), errors.WithMessage(HttpWrongStatusError,
				fmt.Sprintf("wrong status: %d", response.StatusCode()))
		}
		return string(rspBody), nil
	default:
		return "", errors.New("the method is not available")
	}
}

func (d *DataSourceService) queryDBData(item DatabaseSource, params map[string]interface{}) (string, error) {
	var offSet int
	var limit int
	if params["_page"] != nil && params["_size"] != nil {
		page, err := strconv.Atoi(fmt.Sprintf("%v", params["_page"]))
		if err != nil || page == 0 {
			return "", errors.Wrap(err, "page param page sent is not right")
		}
		size, err := strconv.Atoi(fmt.Sprintf("%v", params["_size"]))
		if err != nil || size == 0 {
			return "", errors.Wrap(err, "page param size sent is not right")
		}
		offSet = (page-1) * size
		limit = size
	}

	var orderBy string
	if orderByPtr, ok := params["_orderby"]; ok {
		if _, ok := orderByPtr.(string); ok {
			orderBy = orderByPtr.(string)
		} else {
			orderBy = ""
		}
	}

	engine := dataOwnerMapInstance[item.OrgCode].Engine
	switch strings.ToLower(engine) {
	case "mysql":
		mysqlClient, ok := d.mysqlPool[item.OrgCode]
		if !ok {
			return "", errors.New("can not find sql client by org code")
		}
		instances := make([]map[string]interface{}, 0)
		if item.SqlQueryStr != "" {
			sqlStr := GetSqlStr(item.SqlQueryStr, item.Params, params)
			err := mysqlClient.Query(&instances, sqlStr)
			if err != nil {
				return "", err
			}

		} else {
			whereStr, args := GetMysqlCondition(item.Params, params)

			tableList := strings.Split(item.TableName, ",")
			for _, tName := range tableList {
				tName = strings.Trim(tName, " ")
				if tName == "" {
					continue
				}
				instance := make([]map[string]interface{}, 0)
				condition := sql.Condition{
					Where:		whereStr,
					Args: 		args,
					TableName:  tName,
					Offset:     offSet,
					Limit:      limit,
					OrderBy:    orderBy,
				}
				err := mysqlClient.Find(&condition, &instance)
				if err != nil {
					return "", err
				}
				for _, singleInstance := range instance {
					for k, v := range singleInstance {
						if vByte, ok := v.([]byte); ok {
							singleInstance[k] = string(vByte)
						}
					}
				}
				instances = append(instances, instance...)
			}
		}

		jsonBytes, err := json.Marshal(instances)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	case "mongo":
		mongoClient, ok := d.mongoPool[item.OrgCode]
		if !ok {

			return "", errors.New("can not find mongo client by org code")
		}

		instances := make([]map[string]interface{}, 0)

		tableList := strings.Split(item.TableName, ",")
		for _, tName := range tableList {
			tName = strings.Trim(tName, " ")
			if tName == "" {
				continue
			}
			instance := make([]map[string]interface{}, 0)
			condition := GetMongoCondition(item.Params, params)

			var err error
			if limit != 0 {
				err = mongoClient.FindRange(condition, &instance, offSet, limit, item.TableName)
			} else {
				err = mongoClient.Find(condition, &instance, item.TableName)
			}
			if err != nil {
				return "", err
			}
			instances = append(instances, instance...)
		}

		jsonBytes, err := json.Marshal(instances)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	default:
		return "", errors.New("method is not available")
	}
}

func (d *DataSourceService) queryKVData(item KVDatabaseSource, keys string) (string, error) {
	var keyContents []string
	switch dataOwnerMapInstance[item.OrgCode].Engine {
	case "redis":
		redisClient, ok := d.redisPool[item.OrgCode]
		if !ok {
			return "", errors.New("can not find redis client by org code")
		}
		for _, keyName := range strings.Split(keys, ",") {
			keyName := strings.Trim(keyName, " ")
			if keyName == "" {
				continue
			}
			var err error
			var content string
			action := item.Operation
			if action == "get" {
				content, err = redisClient.GetStringValue(keyName)
			} else if action == "lpop" {
				content, err = redisClient.LPop(keyName)
			} else if action == "rpop" {
				content, err = redisClient.RPop(keyName)
			} else {
				return "", errors.New("the operation is not available")
			}
			if err != nil {
				return "", err
			}
			keyContents = append(keyContents, content)
		}
	default:
		return "", errors.New("the engine is not available")
	}
	keyContentsStr, _ := json.Marshal(keyContents)
	return string(keyContentsStr), nil
}
