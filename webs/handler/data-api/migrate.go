package data_api

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql/mysql"
)

var (
	wrapperInstance = 		&Wrapper{}
	dataApiInstance   = 	&DataApi{}

	tableInstanceMap = map[string]interface{}{
		wrapperName :    	wrapperInstance,
		dataApiName:   		dataApiInstance,
	}
)

func (d *DataApiWeb) Migrate(mysqlClient *mysql.Client, tNames ...string) error {
	if len(tNames) == 0 {
		err := mysqlClient.MigrateByTableName(tableInstanceMap)
		if err != nil {
			return err
		}
	} else {
		filterMap := make(map[string]interface{}, 5)
		if tNames[0] != "all" {
			for _, tName := range tNames {
				if tName == "" {
					continue
				}
				if instance, ok := tableInstanceMap[tName]; ok {
					filterMap[tName] = instance
				} else {
					return errors.New(fmt.Sprintf("receive wrong table name of %s\n", tName))
				}
			}
			err := mysqlClient.MigrateByTableName(filterMap)
			if err != nil {
				return err
			}
		} else {
			err := mysqlClient.MigrateByTableName(tableInstanceMap)
			if err != nil {
				return err
			}
		}

	}
	return nil
}


