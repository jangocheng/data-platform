package oracle

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"xorm.io/core"
)

func GetWhereCondition(where string, args []interface{}) (string, error) {

	needArgsNum := strings.Count(where, "?")
	if needArgsNum > len(args) {
		return "", errors.New(fmt.Sprintf("need %d args, receive %d", needArgsNum, len(args)))
	} else if needArgsNum == 0 {
		return "", nil
	}

	var formatArgs []interface{}

	for _, arg := range args {

		if _, ok := arg.(string); ok {
			formatArgs = append(formatArgs, fmt.Sprintf("'%v'", arg))
		} else {
			formatArgs = append(formatArgs, arg)
		}
	}

	whereFormatStr := strings.Replace(where, "?", "%v", -1)

	whereStr := fmt.Sprintf(whereFormatStr, formatArgs...)

	if whereStr == "" {
		return "", nil
	} else {
		return "where " + whereStr, nil
	}
}


func (c *Client) GetPkName(tableName string) (string, error) {
	dbMeta, err := c.engine.DBMetas()
	if err != nil {
		return "", errors.New(fmt.Sprintf("oracle db meta get fail, error info: %s", err))
	}
	var dbMetaInfo *core.Table
	for _, dbMeta := range dbMeta {
		if dbMeta.Name == tableName {
			dbMetaInfo = dbMeta
			break
		}
	}
	if dbMetaInfo == nil {
		return "", errors.New(fmt.Sprintf("oracle talble can not find, error info: %s", err))
	}
	if len(dbMetaInfo.PrimaryKeys) < 1 {
		return "", errors.New("table has not primary keys")
	}


	return dbMetaInfo.PrimaryKeys[0], nil
}