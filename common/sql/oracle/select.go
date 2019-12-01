package oracle

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql"
)

func (c *Client) Get(condition *sql.Condition, instance interface{}) (map[string]interface{}, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	if pk != nil {
		pkName, err := c.GetPkName(tableName)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		pkValue := ""
		if _, ok := pk.(string); ok {
			pkValue = fmt.Sprintf("'%v'", pk)
		} else {
			pkValue = fmt.Sprintf("%v", pk)
		}
		result, err := c.engine.QueryInterface(fmt.Sprintf("select * from \"%s\" where \"%s\"=%v",
			tableName, pkName, pkValue))
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if len(result) != 0 {
			return result[0], nil
		} else {
			return nil, nil
		}

	} else {
		var sqlStr string
		whereCon, err := GetWhereCondition(where, args)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if whereCon == "" {
			sqlStr = fmt.Sprintf("select * from \"%s\" where rownum<=1", tableName)
		} else {
			sqlStr = fmt.Sprintf("select * from \"%s\" %s and rownum<=1", tableName, whereCon)
		}
		result, err := c.engine.QueryInterface(sqlStr)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if len(result) != 0 {
			return result[0], nil
		} else {
			return nil, nil
		}
	}
}

func (c *Client) Find(condition *sql.Condition) ([]map[string]interface{}, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	limit := condition.Limit
	offset := condition.Offset
	orderBy := condition.OrderBy
	if pk != nil {
		pkName, err := c.GetPkName(tableName)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		pkValue := ""
		if _, ok := pk.(string); ok {
			pkValue = fmt.Sprintf("'%v'", pk)
		} else {
			pkValue = fmt.Sprintf("%v", pk)
		}
		result, err := c.engine.QueryInterface(fmt.Sprintf("select * from \"%s\" where \"%s\"=%v",
			tableName, pkName, pkValue))
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if result == nil {
			result = make([]map[string]interface{}, 0)
		}
		return result, nil
	} else {
		var sqlStr string
		whereCon, err := GetWhereCondition(where, args)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if limit != 0 || offset != 0 {
			if whereCon == "" {
				sqlStr = fmt.Sprintf(
					"select * from (select rownum as no,a.* from \"%s\" a where rownum <= %d) where no > %d",
					tableName, limit+offset, offset)
			} else {
				sqlStr = fmt.Sprintf(
					"select * from (select rownum as no,a.* from \"%s\" a %s and rownum <= %d) where no > %d",
					tableName, whereCon, limit+offset, offset)
			}
		} else {
			sqlStr = fmt.Sprintf(
				"select * from \"%s\" %s", tableName, whereCon)
		}
		if orderBy != "" {
			sqlStr = sqlStr + fmt.Sprintf(" order by \"%s\" desc", orderBy)
		}

		result, err := c.engine.QueryInterface(sqlStr)
		if err != nil {
			return nil, errors.Wrap(err, "oracle instance get error")
		}
		if result == nil {
			result = make([]map[string]interface{}, 0)
		}
		return result, nil
	}
}


func (c *Client) Query(instance *[]map[string]interface{}, sql string) error {
	resultMap, err := c.engine.QueryInterface(sql)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("oracle instance get error"))
	}
	if resultMap == nil {
		resultMap = make([]map[string]interface{}, 0)
	}
	*instance = resultMap
	return nil
}

