package data_derive

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry/mdns"
	"github.com/micro/go-micro/registry/memory"
	"github.com/pkg/errors"
	"platform/common/keeper"
	kmysql "platform/common/keeper/mysql"
	"platform/common/sql"
	"platform/common/sql/mysql"
	"platform/common/utils"
	"platform/common/vm"
	"platform/common/vm/otto"
	"platform/micro-common/log"
	microOtto "platform/micro-common/vm/otto"
	"strings"
	"time"
)

var (

	deriveCodeParamName = "code"
	dataApiParamName = "data"
	dataApiScheme = "http"
	cacheParamName = "cache"

	apiNodes []string

	dataApiWebServiceName = "go.micro.web.data-api"
	dataApiWebSuccessCode int32 = 20000
	dataApiDataPath = "/data/queryApi"
	dataApiParamsPath = "/data/params"


	ErrorStatus int32 = 50000
	SuccessStatus int32 = 20000
	ParamsErrorStatus int32 = 30000
	NotFundErrorStatus int32 = 40000
	InvalidStatus int32 = 10000
	DataDeriveCenterConfigError int32 = 60000
	CallDataApiError int32 = 70000
	DataWrapperJudgeError int32 = 80000
	WrapperExecError int32 = 90000


	wrapperName   		= 	"derive_wrapper"
	dataDeriveName   	= 	"data_derive"
	dataDeriveSetName	= 	"data_derive_set"

	wrapperMapInstance = 		make(map[string]Wrapper)
	dataDeriveMapInstance = 	make(map[string]DataDerive)
	dataDeriveSetMapInstance =  make(map[string]DataDeriveSet)

	tableIndexMap = map[string]string{
		wrapperName:   		"WrapperCode",
		dataDeriveName:   	"DeriveCode",
		dataDeriveSetName:  "DeriveSetCode",
	}
	tableDataInstanceMap = map[string]interface{}{
		wrapperName:      		&wrapperMapInstance,
		dataDeriveName:   		&dataDeriveMapInstance,
		dataDeriveSetName:      &dataDeriveSetMapInstance,
	}
)


type Wrapper struct {
	Id          	int64  `xorm:"pk autoincr(1)"`
	WrapperCode     string `xorm:"varchar(128) notnull unique index"`
	Type            string `xorm:"varchar(32)"`    // before \ after \ log
	JsCode          string `xorm:"longtext"`
	Desc        	string `xorm:"longtext"`
}

type DataDerive struct {
	Id          		int64  `xorm:"pk autoincr(1)"`
	DeriveCode     		string `xorm:"varchar(128) notnull unique index"`
	DeriveName 			string `xorm:"varchar(128) notnull index"`
	ApiCodes			string `xorm:"varchar(256) notnull"`
	Invalid     		int8   `xorm:"tinyint(1) notnull default(0)"`
	Desc        		string `xorm:"longtext"`
	InLineWrapper       string `xorm:"varchar(4096)"`
	AfterWrapperCode    string `xorm:"varchar(512)"`
	LogWrapperCode      string `xorm:"varchar(512)"`
	VerifyParams        string `xorm:"varchar(4096)"`
	CreatedTime 		string `xorm:"datetime created"`
	UpdatedTime 		string `xorm:" datetime updated"`
}

type DataDeriveSet struct {
	Id          		int64  `xorm:"pk autoincr(1)"`
	DeriveSetCode     	string `xorm:"varchar(128) notnull unique index"`
	DeriveCodes         string `xorm:"varchar(2048)"`
	DeriveName 			string `xorm:"varchar(128) notnull index"`
	Desc        		string `xorm:"longtext"`
	LogWrapperCode      string `xorm:"varchar(512)"`
	VerifyParams        string `xorm:"varchar(4096)"`
	CreatedTime 		string `xorm:"datetime created"`
	UpdatedTime 		string `xorm:"datetime updated"`
}



type DataDeriveConfig struct {
	MysqlNode 		string
	MysqlDb   		string
	MysqlUser 		string
	MysqlPs   		string
	LogPath 		string
	Migrate   		string
	Registry    	string
	RegistryAddress string
}


type DataDeriveWeb struct {
	Router  		*mux.Router
	keeper      	keeper.Keeper
	logger      	*log.Logger
	killer          *utils.Killer
	reg				registry.Registry
	ottoClient      vm.ClientVM
}

type DeriveResponse struct {
	Status      		int32
	Message     		string
	Result      		interface{}
	SwiftNumber 		string
	CostTime            int32
}

type ApiDataResponse struct {
	Status      		int32
	Message     		string
	Result      		interface{}
	SwiftNumber 		string
	CostTime            int32
}

type Params struct {
	Params       map[string]string
	Must         []string
}

type ApiParamsResponse struct {
	Status      int32
	Message     string
	Result      Params
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
}

type AfterWrapperMsg struct {
	Fail   	bool
	Params  interface{}
	Result  interface{}
	Item    interface{}
}

type LogWrapperMsg struct {
	Fail   	bool
	Params  map[string]string
	Result 	interface{}
}

type BeforeWrapperMsg map[string]interface{}


func NewDataDeriveWeb(cnf DataDeriveConfig) (*DataDeriveWeb, error) {
	if cnf.MysqlDb == "" || cnf.MysqlPs == "" || cnf.MysqlUser == "" || cnf.MysqlNode == "" {
		return nil, errors.New("config field mysql_db mysql_ps mysql_user mysql_node can not be null")
	}
	var keeperClient keeper.Keeper
	var err error
	var dataDeriveLogger *log.Logger
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
		dataDeriveLogger, err = log.NewLogger()
	} else {
		dataDeriveLogger, err = log.NewLogger(cnf.LogPath)
	}
	if err != nil {
		_ = keeperClient.Close()
		return nil, errors.Wrap(err, "fail to create logger")
	}
	killer := utils.NewKiller()
	d := &DataDeriveWeb{
		keeper:      keeperClient,
		logger:      dataDeriveLogger,
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
	d.reg = reg
	ottoClient, err := otto.NewClientOttoPool(8)
	if err != nil {
		_ = d.Close()
		return nil, errors.Wrap(err, "fail to new a toot client pool")
	}
	err = microOtto.SetFunctionTools(ottoClient)
	if err != nil {
		_ = d.Close()
		return nil, err
	}
	d.ottoClient = ottoClient
	d.startKeepWatch()
	d.startKeepRegService(dataApiWebServiceName)
	d.Router = mux.NewRouter()
	d.RegisterHandler()
	return d, nil
}

func (d *DataDeriveWeb) Close(tNames ...string) error {
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
	return errors.New(errStr)
}


func (d *DataDeriveWeb) startKeepWatch() {
	for tName, instance := range tableDataInstanceMap {
		idxName := tableIndexMap[tName]
		d.keeper.Watch(instance, tName, idxName)
	}
}

func (d *DataDeriveWeb) startKeepRegService(srvName string) {
	go func() {
		for {
			if d.killer.KillNow() {
				break
			}
			apiWebServices, err := d.reg.GetService(srvName)
			if err != nil {
				_ = d.logger.Log("fail to get data api web address from registry")
			}
			if len(apiWebServices) > 0 {
				var newApiNodes []string
				for _, service := range apiWebServices {
					for _, node := range service.Nodes {
						newApiNodes = append(newApiNodes, node.Address)
					}
				}
				apiNodes = newApiNodes
			} else {
				apiNodes = []string{}
			}

			time.Sleep(8 * time.Second)
		}
	}()
}

func (d *DataDeriveWeb) PP() {
	for {
		if d.killer.KillNow() {
			break
		}
		time.Sleep(time.Second * 6)
		fmt.Println("keeper", d.keeper)
		fmt.Println("js_wrapper", wrapperMapInstance)
		fmt.Println("data_derive", dataDeriveMapInstance)
	}
}
