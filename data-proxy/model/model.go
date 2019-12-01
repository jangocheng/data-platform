package model

import (
	"fmt"
	"platform/common/redis"
	plog "platform/micro-common/log"
	"strings"
)

var (
	SuccessCode 		int32 = 20000
	ParamsErrorCode 	int32 = 40000
	ServerErrorCode     int32 = 50000
	NoResultCode      	int32 = 80011
	QueryErrorCode      int32 = 80000
	InvalidApiCode		int32 = 90000
	NoPermissionCode    int32 = 90001
	QueryLimitCode      int32 = 90002

	PermissionUrlMap = map[string]string{
		"/ice_wall/api/get_list": "",
		"/ice_wall/api/get_data": "",
		"/ice_wall/api/submit": "",
	}
	GlobalConf  = new(Config)

	RedisClient *redis.ClientRedis
)

type DataPlatformResponse struct {
	Status      		int32
	Message     		string
	Result      		interface{}
	SwiftNumber 		string
	CostTime            int32
}

type Config struct {
	RedisNode       		string
	RedisPs         		string
	RedisDb					int

	BairongHost     		string
	BairongLoginPath 		string
	BairongUserName     	string
	BairongPassword     	string

	BairongBizHost          string
	BairongBizApiPath       string
	BairongBizDerivePath    string

	PlatformHost            string
	PlatformApiPath         string
	PlatformDeriveSetPath   string
	PlatformDerivePath      string

	StrategyHost            string
	StrategyListUrl         string
	StrategyStartUrl        string

	PostLoanHost            string
	PostLoanListUrl         string
	PostLoanStartUrl        string

	NeedCacheId             []string

	LogPath                 string

	DefaultApiCode          string

	Debug                   bool
}

func (c *Config) NeedCache(flowId string) bool {
	for _, _id := range c.NeedCacheId {
		if _id == flowId {
			return true
		}
	}
	return false
}

type Runtime struct {
	Status              int32
	RequestTime         string
	CostTime   			int32
	FormData            interface{}
	Message             string
}

type Response struct {
	Code   		int32
	Message     string
	Result      interface{}
	SwiftNumber string
}

var Logger *plog.Logger

func InitModel() {
	var err error
	if GlobalConf.LogPath != "" {
		Logger, err = plog.NewLogger(GlobalConf.LogPath)
	} else {
		Logger, err = plog.NewLogger()
	}
	if err != nil {
		panic(fmt.Sprintf("fail to create logger error info: %s", err.Error()))
	}


	if GlobalConf.RedisNode != "" {
		cacheCnf := redis.ConnectConfig{
			IpPortList: strings.Split(GlobalConf.RedisNode, ","),
			Password:   GlobalConf.RedisPs,
			PoolSize:   32,
			DB:         GlobalConf.RedisDb,
		}
		RedisClient , err = redis.NewRedisClient(cacheCnf)
		if err != nil {
			panic("fail to start redis client")
		}
	}
}
