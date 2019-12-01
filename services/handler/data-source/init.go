package data_source

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"platform/common/keeper"
	kmysql "platform/common/keeper/mysql"
	"platform/common/mongo"
	"platform/common/redis"
	"platform/common/sql"
	"platform/common/sql/mysql"
	"platform/common/utils"
	"platform/micro-common/log"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorStatus int32 = 5000
	SuccessStatus int32 = 2000
	ParamsErrorStatus int32 = 3000
	NotFundErrorStatus int32 = 4000
	InvalidStatus int32 = 1000
	HttpWrongStatus int32 = 6000



	dataOwnerName    = "data_owner"
	fileDataSource   = "file_data_source"
	httpDataSource   = "http_data_source"
	databaseSource   = "database_source"
	kvDatabaseSource = "kv_database_source"

	dataOwnerMapInstance = 			make(map[string]DataOwner)
	fileDataSourceMapInstance = 	make(map[string]FileDataSource)
	httpDataSourceMapInstance = 	make(map[string]HttpDataSource)
	databaseSourceMapInstance = 	make(map[string]DatabaseSource)
	kvDatabaseSourceMapInstance = 	make(map[string]KVDatabaseSource)

	tableIndexMap = map[string]string{
		dataOwnerName:    "OrgCode",
		fileDataSource:   "DsCode",
		httpDataSource:   "DsCode",
		databaseSource:   "DsCode",
		kvDatabaseSource: "DsCode",
	}

	tableDataInstanceMap = map[string]interface{}{
		dataOwnerName:      	&dataOwnerMapInstance,
		fileDataSource:   		&fileDataSourceMapInstance,
		httpDataSource:			&httpDataSourceMapInstance,
		databaseSource: 		&databaseSourceMapInstance,
		kvDatabaseSource:		&kvDatabaseSourceMapInstance,
	}
)

type logFormat struct {
	Status              int32
	Path       			string
	RequestTime         string
	CostTime   			int32
	ParentSwiftNumber 	string
	FormData    		interface{}
	Result   			interface{}
	Message             string
}


type DataOwner struct {
	Id          int64  `xorm:"pk autoincr(1)"`
	OrgCode     string `xorm:"varchar(128) notnull unique index"`
	OrgName     string `xorm:"varchar(128) notnull index"`
	SelfOwn     int8   `xorm:"tinyint(1) notnull"`
	Invalid     int8   `xorm:"tinyint(1) notnull default(0)"`
	Addr        string `xorm:"varchar(128)"`
	User        string `xorm:"varchar(32)"`
	PassWord    string `xorm:"varchar(64)"`
	Db          string `xorm:"varchar(32)"`
	Engine      string `xorm:"varchar(16)"` // file / http / mysql / mongo / oracle / redis
	Desc        string `xorm:"longtext"`
	CreatedTime string `xorm:"datetime created"`
	UpdatedTime string `xorm:" datetime updated"`
}

type FileDataSource struct {
	Id          int64  `xorm:"pk autoincr(1)"`
	OrgCode     string `xorm:"varchar(128) notnull index"`
	DsCode      string `xorm:"varchar(128) notnull unique index"`
	DsName      string `xorm:"varchar(128) notnull index"`
	Invalid     int8   `xorm:"tinyint(1) notnull default(0)"`
	FileFolder  string `xorm:"varchar(256) notnull default('/')"`
	Desc        string `xorm:"longtext"`
	CreatedTime string `xorm:"datetime created"`
	UpdatedTime string `xorm:" datetime updated"`
}

type HttpDataSource struct {
	Id          int64  `xorm:"pk autoincr(1)"`
	OrgCode     string `xorm:"varchar(128) notnull index"`
	DsCode      string `xorm:"varchar(128) notnull unique index"`
	DsName      string `xorm:"varchar(128) notnull index"`
	Invalid     int8   `xorm:"tinyint(1) notnull default(0)"`
	Desc        string `xorm:"longtext"`
	Method      string `xorm:"varchar(8)"`
	SuccessCode string `xorm:"varchar(8)"`
	Path        string `xorm:"varchar(256) notnull default('/')"`
	Params      string `xorm:"varchar(256)"`
	ContentType string `xorm:"varchar(32)"` // json/form
	CreatedTime string `xorm:"datetime created"`
	UpdatedTime string `xorm:" datetime updated"`
}

type DatabaseSource struct {
	Id          	int64  `xorm:"pk autoincr(1)"`
	OrgCode     	string `xorm:"varchar(128) notnull index"`
	DsCode      	string `xorm:"varchar(128) notnull unique index"`
	DsName      	string `xorm:"varchar(128) notnull index"`
	Invalid     	int8   `xorm:"tinyint(1) notnull default(0)"`
	Desc        	string `xorm:"longtext"`
	TableName   	string `xorm:"varchar(128) notnull"`
	SqlQueryStr 	string `xorm:"varchar(512)"`
	Params      	string `xorm:"varchar(256)"`
	CreatedTime 	string `xorm:"datetime created"`
	UpdatedTime 	string `xorm:" datetime updated"`
}

type KVDatabaseSource struct {
	Id          int64  `xorm:"pk autoincr(1)"`
	OrgCode     string `xorm:"varchar(128) notnull index"`
	DsCode      string `xorm:"varchar(128) notnull unique index"`
	DsName      string `xorm:"varchar(128) notnull index"`
	Invalid     int8   `xorm:"tinyint(1) notnull default(0)"`
	Desc        string `xorm:"longtext"`
	Operation   string `xorm:"varchar(32) notnull"`
	CreatedTime string `xorm:"datetime created"`
	UpdatedTime string `xorm:" datetime updated"`
}

type DataSourceConfig struct {
	MysqlNode string
	MysqlDb   string
	MysqlUser string
	MysqlPs   string
	LogPath string
	Migrate   string
}

type DataSourceService struct {
	keeper      	keeper.Keeper
	logger      	*log.Logger
	mysqlPool   	map[string]*mysql.Client
	redisPool   	map[string]*redis.ClientRedis
	mongoPool       map[string]*mongo.ClientMongo
	killer          *utils.Killer
}

func NewDataSourceService(cnf *DataSourceConfig) (*DataSourceService, error) {
	if cnf.MysqlDb == "" || cnf.MysqlPs == "" || cnf.MysqlUser == "" || cnf.MysqlNode == "" {
		return nil, errors.New("config field mysql_db mysql_ps mysql_user mysql_node can not be null")
	}

	var keeperClient keeper.Keeper
	var err error
	var dataSourceLogger *log.Logger
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
		dataSourceLogger, err = log.NewLogger()
	} else {
		dataSourceLogger, err = log.NewLogger(cnf.LogPath)
	}
	if err != nil {
		_ = keeperClient.Close()
		return nil, errors.Wrap(err, "fail to create logger")
	}
	killer := utils.NewKiller()
	d := &DataSourceService{
		keeper:      keeperClient,
		logger:      dataSourceLogger,
		mysqlPool:   make(map[string]*mysql.Client),
		redisPool:   make(map[string]*redis.ClientRedis),
		mongoPool:   make(map[string]*mongo.ClientMongo),
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
	startKeepWatch(keeperClient)
	go d.KeepDBPools()
	return d, nil
}

func (d *DataSourceService) KeepDBPools(tNames ...string) {
	orgItemHashMap := make(map[string]string)
	for {
		if d.killer.KillNow() {
			return
		}
		for orgCode, orgItem := range dataOwnerMapInstance {
			if d.killer.KillNow() {
				return
			}
			var errStr string
			if orgItem.Engine != "" && orgItem.Db != "" && orgItem.PassWord != "" &&
				orgItem.Addr != ""{
				switch orgItem.Engine {
				case "mysql":
					jsonBytes, _ := json.Marshal(orgItem)
					itemHash := utils.GetMd5(string(jsonBytes))
					if _, ok := d.mysqlPool[orgCode]; ok {
						if _itemHash, ok := orgItemHashMap[orgCode]; ok {
							if itemHash == _itemHash {
								continue
							}
						}
					}
					conf := sql.ConfigSql{
						NodeList: strings.Split(orgItem.Addr, ","),
						User:  orgItem.User,
						Passwd: orgItem.PassWord,
						DB: orgItem.Db,
						MaxIdleConns: 8,
						MaxOpenConns: 32,
					}
					mysqlClient, err := mysql.NewClient(conf)
					if err != nil {
						errStr = err.Error()
						break
					}
					if client, ok := d.mysqlPool[orgCode]; ok {
						if client != nil {
							_ = client.Close()
						}
					}
					orgItemHashMap[orgCode] = itemHash
					d.mysqlPool[orgCode] = mysqlClient
				case "redis":
					jsonBytes, _ := json.Marshal(orgItem)
					itemHash := utils.GetMd5(string(jsonBytes))
					if _, ok := d.redisPool[orgCode]; ok {
						if _itemHash, ok := orgItemHashMap[orgCode]; ok {
							if itemHash == _itemHash {
								continue
							}
						}
					}
					dbInt, err := strconv.Atoi(orgItem.Db)
					if err != nil {
						errStr = err.Error()
						break
					}
					conf := redis.ConnectConfig{
						IpPortList:   strings.Split(orgItem.Addr, ","),
						Password:     orgItem.PassWord,
						PoolSize:     16,
						MinIdexCon:   4,
						DB:           dbInt,
					}
					redisClient, err := redis.NewRedisClient(conf)
					if err != nil {
						errStr = err.Error()
						break
					}
					if client, ok := d.redisPool[orgCode]; ok {
						if client != nil {
							_ = client.Close()
						}
					}
					orgItemHashMap[orgCode] = itemHash
					d.redisPool[orgCode] = redisClient
				case "mongo":
					jsonBytes, _ := json.Marshal(orgItem)
					itemHash := utils.GetMd5(string(jsonBytes))
					if _, ok := d.mongoPool[orgCode]; ok {
						if _itemHash, ok := orgItemHashMap[orgCode]; ok {
							if itemHash == _itemHash {
								continue
							}
						}
					}

					conf := mongo.ConnectConfig{
						IpPortList: 	strings.Split(orgItem.Addr, ","),
						UserName:  		orgItem.User,
						Password: 		orgItem.PassWord,
						DB: 			orgItem.Db,
					}
					mongoClient, err := mongo.NewClient(conf)
					if err != nil {
						errStr = err.Error()
						break
					}
					if client, ok := d.mongoPool[orgCode]; ok {
						if client != nil {
							client.Close()
						}
					}
					orgItemHashMap[orgCode] = itemHash
					d.mongoPool[orgCode] = mongoClient
				default:
					errStr = fmt.Sprintf("new database client occur error, org_code: %s engine: %s not support",
						orgCode, orgItem.Engine)
				}
			}
			if errStr != "" {
				_ = d.logger.Log(fmt.Sprintf("new database client occur error, org_code: %s err: %s",
					orgCode, errStr), "WARN")
			}
		}


		for orgCode, client := range d.mongoPool {
			if _, ok := dataOwnerMapInstance[orgCode]; !ok {
				client.Close()
				delete(d.mongoPool, orgCode)
				delete(orgItemHashMap, orgCode)
			}
		}
		for orgCode, client := range d.redisPool {
			if _, ok := dataOwnerMapInstance[orgCode]; !ok {
				_ = client.Close()
				delete(d.redisPool, orgCode)
				delete(orgItemHashMap, orgCode)
			}
		}
		for orgCode, client := range d.mysqlPool {
			if _, ok := dataOwnerMapInstance[orgCode]; !ok {
				_ = client.Close()
				delete(d.mysqlPool, orgCode)
				delete(orgItemHashMap, orgCode)
			}
		}
		time.Sleep(time.Second * 16)
	}
}

func (d *DataSourceService) PP() {
	for {
		if d.killer.KillNow() {
			break
		}
		time.Sleep(time.Second * 16)
		fmt.Println("redis_pool", d.redisPool)
		fmt.Println("sql_pool", d.mysqlPool)
		fmt.Println("mongo_pool", d.mongoPool)
		fmt.Println("keeper", d.keeper)
		fmt.Println("owner", dataOwnerMapInstance)
		fmt.Println("file", fileDataSourceMapInstance)
		fmt.Println("http", httpDataSourceMapInstance)
		fmt.Println("db", databaseSourceMapInstance)
		fmt.Println("kv", kvDatabaseSourceMapInstance)
	}
}

func (d *DataSourceService) Close(tNames ...string) error {
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
	if d.mysqlPool != nil {
		for _, sqlClient := range d.mysqlPool {
			errSql := sqlClient.Close()
			if errSql != nil {
				errStr = errStr + errSql.Error()
			}
		}
	}
	if d.redisPool != nil {
		for _, redisClient := range d.redisPool {
			errRedis := redisClient.Close()
			if errRedis != nil {
				errStr = errStr + errRedis.Error()
			}
		}
	}
	if d.mongoPool != nil {
		for _, mongoClient := range d.mongoPool {
			mongoClient.Close()
		}
	}
	return errors.New(errStr)
}


func startKeepWatch(keeperClient keeper.Keeper) {
	for tName, instance := range tableDataInstanceMap {
		idxName := tableIndexMap[tName]
		keeperClient.Watch(instance, tName, idxName)
	}
}