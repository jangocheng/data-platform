package sql

type Condition struct {
	TableName      	string
	Where          	string
	Pk             	interface{}
	Args           	[]interface{}
	Offset         	int
	Limit      		int
	OrderBy          string
}

type ConfigSql struct {
	NodeList 		[]string
	User 			string
	Passwd			string
	DB     			string
	MaxIdleConns 	int
	MaxOpenConns    int
}
