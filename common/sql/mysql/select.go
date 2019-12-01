package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql"
)

func (c *Client) Get(condition *sql.Condition, instance interface{}) (bool, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	var result bool
	var err error
	if pk != nil {
		result, err = c.engine.Table(tableName).ID(pk).Get(instance)
	} else {
		result, err = c.engine.Table(tableName).Where(where, args...).Get(instance)
	}
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("mysql instance get error"))
	}
	return result, nil
}

func (c *Client) Find(condition *sql.Condition, instance interface{}, condiBean ...interface{}) error {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	limit := condition.Limit
	off := condition.Offset
	orderBy := condition.OrderBy

	var err error
	if pk != nil {
		err = c.engine.Table(tableName).ID(pk).Desc(orderBy).Find(instance, condiBean...)
	} else {
		err = c.engine.Table(tableName).Where(where, args...).Limit(limit, off).Find(instance, condiBean...)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("mysql instance find error"))
	}
	return nil
}

func (c *Client) Exist(condition *sql.Condition, instances interface{}) (bool, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args

	var exist bool
	var err error
	if pk != nil {
		exist, err = c.engine.Table(tableName).ID(pk).Exist(instances)
	} else {
		exist, err = c.engine.Table(tableName).Where(where, args...).Exist()
	}

	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("mysql instance find error"))
	}
	return exist, nil
}

func (c *Client) Total(condition *sql.Condition) (int64, error) {
	tableName := condition.TableName
	where := condition.Where
	args := condition.Args

	var total int64
	var err error
	total, err = c.engine.Table(tableName).Where(where, args...).Count()
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("mysql instance find error"))
	}

	return total, nil
}

func (c *Client) Query(instance *[]map[string]interface{}, sql string) error {
	resultMap, err := c.engine.QueryInterface(sql)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("mysql instance get error"))
	}
	if resultMap == nil {
		resultMap = make([]map[string]interface{}, 0)
	}
	*instance = resultMap
	return nil
}

