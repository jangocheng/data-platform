package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql"
)

func (c *Client) Update(condition *sql.Condition, instance interface{}, condiBean ...interface{}) (int64, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	limit := condition.Limit
	off := condition.Offset

	var affected int64
	var err error
	if pk != nil {
		affected, err = c.engine.Table(tableName).ID(pk).Update(instance, condiBean...)
	} else {
		affected, err = c.engine.Table(tableName).Where(where, args...).Limit(limit, off).Update(instance, condiBean...)
	}
	if err != nil {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance update error"))
	} else if affected == 0 {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance update error, affected row %d", affected))
	}
	return affected, nil
}

func (c *Client) UpdateCols(condition *sql.Condition, colsUp map[string]interface{}, condiBean ...interface{}) (int64, error) {
	tableName := condition.TableName
	where := condition.Where
	pk := condition.Pk
	args := condition.Args
	limit := condition.Limit
	off := condition.Offset

	colsSlice := make([]string, len(colsUp))
	for k := range colsUp {
		colsSlice = append(colsSlice, k)
	}

	var affected int64
	var err error
	if pk != nil {
		affected, err = c.engine.Table(tableName).ID(pk).Cols(colsSlice...).Update(colsUp, condiBean...)
	} else {
		affected, err = c.engine.Table(tableName).Where(where, args...).Limit(limit, off).Cols(colsSlice...).Update(colsUp, condiBean...)
	}
	if err != nil {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance update error"))
	} else if affected == 0 {
		return affected, errors.Wrap(err, fmt.Sprintf("mysql instance update error, affected row %d", affected))
	}
	return affected, nil
}
