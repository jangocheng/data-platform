package data_derive

import (
	"fmt"
	"github.com/pkg/errors"
	"platform/common/sql/mysql"
)

var (
	wrapperInstance = 		&Wrapper{}
	dataDeriveInstance   = 	&DataDerive{}
	dataDeriveSetInstance = &DataDeriveSet{}

	tableInstanceMap = map[string]interface{}{
		wrapperName :    		wrapperInstance,
		dataDeriveName:   		dataDeriveInstance,
		dataDeriveSetName:      dataDeriveSetInstance,
	}
)

func (d *DataDeriveWeb) Migrate(mysqlClient *mysql.Client, tNames ...string) error {
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


