package mysql

import "github.com/pkg/errors"

func (c *Client) Migrate (beans ...interface{}) error {
	err := c.engine.Sync2(beans...)
	if err != nil {
		return errors.Wrap(err, "fail to migrate table")
	}
	return nil
}


func (c *Client) MigrateByTableName (beans map[string]interface{}) error {
	errStr := "fail to migrate table by name"
	session := c.engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return errors.Wrap(err, errStr)
	}
	for tableName, tableInstance := range beans {
		err = session.Table(tableName).Sync2(tableInstance)
		if err != nil {
			_ = session.Rollback()
			return errors.Wrap(err, errStr)
		}
	}
	err = session.Commit()
	if err != nil {
		_ = session.Rollback()
		return errors.Wrap(err, "fail to migrate table by name")
	}
	return nil

}