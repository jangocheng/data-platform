package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"platform/common/keeper"
	k "platform/common/utils"
)

type Keeper struct {
	engine 		*xorm.EngineGroup
	killer   	*k.Killer
	internal	int32
}


type Conf struct {
	NodeList 	[]string
	User 		string
	Passwd		string
	DB     		string
	Internal    int32
}

func NewKeeper(conf Conf) (keeper.Keeper, error) {
	dnsList := make([]string, 0)
	for _, node := range conf.NodeList {
		dnsString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.User, conf.Passwd, node, conf.DB)
		dnsList = append(dnsList, dnsString)
	}
	engineGroup, err := xorm.NewEngineGroup("mysql", dnsList)
	if err != nil {
		return nil, errors.Wrap(err, "fail to new mysql keeper")
	}
	engineGroup.SetMaxOpenConns(4)
	killer := k.NewKiller()
	return &Keeper{engine: engineGroup, killer: killer, internal: conf.Internal}, nil
}

func (k *Keeper) Close() error {
	if k.engine != nil {
		return k.engine.Close()
	}
	return nil
}



