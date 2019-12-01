package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql"
)


func (c *Client) Delete(condition *sql.Condition, instance interface{}) (int64, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	off := condition.Offset
	limit := condition.Limit
	var affected int64
	var err error
	if pk != nil {
		affected, err = c.engine.Table(tableName).ID(pk).Unscoped().Delete(instance)
	} else {
		affected, err = c.engine.Table(tableName).Where(where, args...).Limit(limit, off).Unscoped().Delete(instance)
	}
	if err != nil {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance delete error"))
	} else if affected == 0 {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance delete error, affected row %d", affected))
	}
	return affected, nil
}

