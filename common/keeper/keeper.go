package keeper

type Keeper interface {
	Watch(instance interface{}, tableName string, idxName string)
	Close() error
}




