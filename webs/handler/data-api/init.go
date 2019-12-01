package data_api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry/mdns"
	"github.com/micro/go-micro/registry/memory"
	"github.com/pkg/errors"
	"math/rand"
	"platform/common/keeper"
	kmysql "platform/common/keeper/mysql"
	"platform/common/redis"
	"platform/common/sql"
	"platform/common/sql/mysql"
	"platform/common/utils"
	"platform/common/vm"
	"platform/common/vm/otto"
	"platform/micro-common/log"
	microOtto "platform/micro-common/vm/otto"
	proto "platform/proto/service/data-source"
	"strings"
	"time"
)

var (

	apiCodeParamName = "code"
	dataSourceParamName = "data"
	cacheParamName = "cache"

	dataSourceRpcServiceName = "go.micro.srv.data-source"
	dataSourceRpcServiceSuccessCode int32 = 2000

	ErrorStatus int32 = 50000
	CacheError  int32 = 50020
	SuccessStatus int32 = 20000
	ParamsErrorStatus int32 = 30000
	NotFundErrorStatus int32 = 40000
	InvalidStatus int32 = 10000
	DataApiCenterConfigError int32 = 60000
	DataSourceError int32 = 70000
	DataWrapperJudgeError int32 = 80000
	WrapperExecError int32 = 90000


	wrapperName   = "api_wrapper"
	dataApiName   = "data_api"

	wrapperMapInstance = 	make(map[string]Wrapper)
	dataApiMapInstance = 	make(map[string]DataApi)

	tableIndexMap = map[string]string{
		wrapperName:   "WrapperCode",
		dataApiName:   "ApiCode",
	}
	tableDataInstanceMap = map[string]interface{}{
		wrapperName:      	&wrapperMapInstance,
		dataApiName:   		&dataApiMapInstance,
	}
)

type Cache struct {
	UseCache bool
	Hit      bool
}

type LogFormat struct {
	Status              int32
	Path       			string
	RequestTime         string
	CostTime   			int32
	SwiftNumber         string
	ParentSwiftNumber   string
	RawData             string
	FormData            interface{}
	Result   			interface{}
	LogWrapperTime      int32
	Message             string
	Cache               Cache
}

type AfterWrapperMsg struct {
	Fail   	bool
	Params  interface{}
	Result 	interface{}
	Item    interface{}
}

type LogWrapperMsg struct {
	Fail   	bool
	Params  map[string]string
	Result 	interface{}
}

type BeforeWrapperMsg map[string]interface{}

type Wrapper struct {
	Id          	int64  `xorm:"pk autoincr(1)"`
	WrapperCode     string `xorm:"varchar(128) notnull unique index"`
	Type            string `xorm:"varchar(32)"`    // before \ after \ log
	JsCode          string `xorm:"longtext"`
	Desc        	string `xorm:"longtext"`
}

type DataApi struct {
	Id          		int64  `xorm:"pk autoincr(1)"`
	ApiCode     		string `xorm:"varchar(128) notnull unique index"`
	ApiName 			string `xorm:"varchar(128) notnull index"`
	DsCode				string `xorm:"varchar(128) notnull index"`
	DsType              string `xorm:"varchar(128) notnull"`
	Invalid     		int8   `xorm:"tinyint(1) notnull default(0)"`
	Desc        		string `xorm:"longtext"`
	BeforeWrapperCode   string `xorm:"varchar(512)"`
	AfterWrapperCode    string `xorm:"varchar(512)"`
	LogWrapperCode      string `xorm:"varchar(512)"`
	VerifyParams        string `xorm:"varchar(4096)"`
	CreatedTime 		string `xorm:"datetime created"`
	UpdatedTime 		string `xorm:" datetime updated"`
}

type DataApiConfig struct {
	MysqlNode 		string
	MysqlDb   		string
	MysqlUser 		string
	MysqlPs   		string
	LogPath 		string
	Migrate   		string
	Registry    	string
	RegistryAddress string
	CacheNode       string
	CachePs         string
	CacheDb			int
	CacheExpire     int
}

type DataApiWeb struct {
	Router  		*mux.Router
	keeper      	keeper.Keeper
	logger      	*log.Logger
	killer          *utils.Killer
	dsClient        proto.DataSourceService
	ottoClient      vm.ClientVM
	cacheClient     *redis.ClientRedis
	cacheExpire     int
}

type ApiResponse struct {
	Status      		int32
	Message     		string
	Result      		interface{}
	SwiftNumber 		string
	CostTime            int32
}

type RpcResultRsp struct {
	Fail     bool
	Result   string
}


func NewDataApiWeb(cnf DataApiConfig) (*DataApiWeb, error) {
	if cnf.MysqlDb == "" || cnf.MysqlPs == "" || cnf.MysqlUser == "" || cnf.MysqlNode == "" {
		return nil, errors.New("config field mysql_db mysql_ps mysql_user mysql_node can not be null")
	}
	var keeperClient keeper.Keeper
	var err error
	var dataApiLogger *log.Logger
	var reg registry.Registry
	mysqlNodeList := strings.Split(cnf.MysqlNode, ",")
	keeperConfig := kmysql.Conf{
		NodeList: mysqlNodeList,
		User:     cnf.MysqlUser,
		Passwd:   cnf.MysqlPs,
		DB:       cnf.MysqlDb,
		Internal: 8,
	}
	if keeperClient, err = kmysql.NewKeeper(keeperConfig); err != nil {
		return nil, errors.Wrap(err, "fail to new keeper client")
	}
	if cnf.LogPath == "" {
		dataApiLogger, err = log.NewLogger()
	} else {
		dataApiLogger, err = log.NewLogger(cnf.LogPath)
	}
	if err != nil {
		_ = keeperClient.Close()
		return nil, errors.Wrap(err, "fail to create logger")
	}
	killer := utils.NewKiller()
	d := &DataApiWeb{
		keeper: keeperClient,
		logger: dataApiLogger,
		killer: killer,
	}
	if cnf.Migrate != "" {
		var mysqlMigClient *mysql.Client
		mysqlCnfConfig := sql.ConfigSql{
			NodeList:     mysqlNodeList,
			User:         cnf.MysqlUser,
			Passwd:       cnf.MysqlPs,
			DB:           cnf.MysqlDb,
			MaxIdleConns: 2,
			MaxOpenConns: 4,
		}
		if mysqlMigClient, err = mysql.NewClient(mysqlCnfConfig); err != nil {
			return nil, errors.Wrap(err, "fail to new mysql migrate client")
		}
		defer mysqlMigClient.Close()
		_ = d.Migrate(mysqlMigClient, strings.Split(cnf.Migrate, ",")...)
	}

	if cnf.Registry == "" {
		reg = mdns.NewRegistry()
	} else if cnf.Registry == "memory" {
		reg = memory.NewRegistry()
	} else if cnf.Registry == "consul" {
		reg = consul.NewRegistry(func(op *registry.Options) {
			op.Addrs = strings.Split(cnf.RegistryAddress, ",")
		})
	} else {
		_ = d.Close()
		return nil, errors.New(fmt.Sprintf("receive wrong registry: %s", cnf.Registry))
	}
	if reg == nil {
		_ = d.Close()
		return nil, errors.New("fail to new registry")
	}
	//service := micro.NewService(micro.Registry(reg))
	//if service == nil {
	//	_ = d.Close()
	//	return nil, errors.New("fail to new rpc service")
	//}
	//service.Init()
	rpcClient := client.NewClient(client.Registry(reg))
	d.dsClient = proto.NewDataSourceService(dataSourceRpcServiceName, rpcClient)
	ottoClient, err := otto.NewClientOttoPool(8)
	if err != nil {
		_ = d.Close()
		return nil, errors.Wrap(err, "fail to new a toot client pool")
	}
	d.ottoClient = ottoClient
	err = microOtto.SetFunctionTools(ottoClient)
	if err != nil {
		_ = d.Close()
		return nil, err
	}
	if cnf.CacheNode != "" {
		cacheCnf := redis.ConnectConfig{
			IpPortList: strings.Split(cnf.CacheNode, ","),
			Password:   cnf.CachePs,
			PoolSize:   32,
			DB:         cnf.CacheDb,
		}
		cacheClient, err := redis.NewRedisClient(cacheCnf)
		if err != nil {
			_ = d.Close()
			return nil, err
		}
		d.cacheClient = cacheClient
		if cnf.CacheExpire != 0 {
			d.cacheExpire = cnf.CacheExpire
		} else {
			d.cacheExpire = 3600
		}
	}
	startKeepWatch(keeperClient)
	d.Router = mux.NewRouter()
	d.RegisterHandler()
	rand.Seed(time.Now().Unix())
	return d, nil
}

func (d *DataApiWeb) Close(tNames ...string) error {
	var errStr string
	if d.keeper != nil {
		errKeeper := d.keeper.Close()
		if errKeeper != nil {
			errStr = errStr + errKeeper.Error()
		}
	}
	if d.logger != nil {
		d.logger.Close()
	}
	if d.ottoClient != nil {
		d.ottoClient.Close()
	}
	if d.cacheClient != nil {
		_ = d.cacheClient.Close()
	}
	return errors.New(errStr)
}

func startKeepWatch(keeperClient keeper.Keeper) {
	for tName, instance := range tableDataInstanceMap {
		idxName := tableIndexMap[tName]
		keeperClient.Watch(instance, tName, idxName)
	}
}

func (d *DataApiWeb) PP() {
	for {
		if d.killer.KillNow() {
			break
		}
		time.Sleep(time.Second * 6)
		fmt.Println("keeper", d.keeper)
		fmt.Println("api_wrapper", wrapperMapInstance)
		fmt.Println("data_api", dataApiMapInstance)
	}
}
