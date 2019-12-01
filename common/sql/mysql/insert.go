package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql"
)

func (c *Client) Insert(condition *sql.Condition, instance ...interface{}) (int64, error) {
	tableName := condition.TableName

	fieldsVerify, failField := mustFieldNotSetVerify(instance)
	if !fieldsVerify {
		return 0, errors.New(fmt.Sprintf("mysql field verify fail, not null field %s has not been set", failField))
	}

	var affected int64
	var err error
	affected, err = c.engine.Table(tableName).Insert(instance...)
	if err != nil {
		return affected, errors.Wrap(err, "mysql instance insert error")
	} else if affected == 0 {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance insert error, affected row %d", affected))
	}
	return affected, nil
}
