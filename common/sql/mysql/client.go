package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"platform/common/sql"
)


type Client struct {
	engine *xorm.EngineGroup
}


func NewClient(conf sql.ConfigSql) (*Client, error) {
	dnsList := make([]string, 0)
	for _, node := range conf.NodeList {
		dnsString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.User, conf.Passwd, node, conf.DB)
		dnsList = append(dnsList, dnsString)
	}
	engineGroup, err := xorm.NewEngineGroup("mysql", dnsList)

	if err != nil {
		return nil, errors.Wrap(err, "fail to new mysql client")
	}
	engineGroup.SetMaxIdleConns(conf.MaxIdleConns)
	engineGroup.SetMaxOpenConns(conf.MaxOpenConns)
	return &Client{engine: engineGroup}, nil
}

func (c *Client) Close() error {
	if c.engine != nil {
		return c.engine.Close()
	}
	return nil
}
